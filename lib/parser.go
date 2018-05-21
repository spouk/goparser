package lib

import (
	"log"
	"io"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"fmt"
	"time"
	"crypto/tls"
	"net/http"
	"net/url"
	"io/ioutil"

	"strings"
	"os"
)

const (
	LOG_PARSER_PREFIX = "[parser-log] "
	LOG_PARSER_FLAGS  = log.Ldate | log.Ltime | log.Lshortfile
)

type Parser struct {
	sync.WaitGroup
	sync.RWMutex
	db                 *SqliteDBS
	log                *log.Logger
	StockSearchResult  []SearchResult
	StockSearchRequest []SearchRequest
	ChanStockRequest   chan SearchRequest
	ChanStockResult    chan SearchResult
	ChanCommand        chan string
	config             *Config
	requestparser      *RequestParser
	transliter         *Transliter
}

func NewParser(config string, logout io.Writer) *Parser {
	//создаю инстанс нового парсера
	p := &Parser{
		log:              log.New(logout, LOG_PARSER_PREFIX, LOG_PARSER_FLAGS),
		ChanCommand:      make(chan string),
		ChanStockRequest: make(chan SearchRequest),
		ChanStockResult:  make(chan SearchResult),
		transliter:       NewTransliter(),
	}
	//парсю конфиг
	p.config = NewConfig(config, logout)
	//создаю/открываю базу данных для манипуляций
	p.db = NewSqliteDBS(p.config.DB, logout)
	//парсю файл с запросами
	p.requestparser = NewRequestParser(p.config, p)
	return p
}
func (p *Parser) Run(countWorkerRequest int, countWorkerFile int) {
	//менеджер по работе с воркерами
	//var totalWorker = countWorkerFile + countWorkerRequest
	//запуск WorkerRequest
	for i := 0; i < countWorkerRequest; i ++ {
		go p.workerRequest(fmt.Sprintf("WORKER_REQUEST#%d", i))
	}
	p.Add(countWorkerRequest)
	p.Wait()

	for i := 0; i < countWorkerFile; i ++ {
		go p.workerFile(fmt.Sprintf("WORKER_FILE#%d", i))
	}
	p.Add(countWorkerFile)
	//запуск WorkerFile
	p.Wait()

	//for totalWorker > 0 {
	//	select {
	//	case result := <-p.ChanStockResult:
	//		wr := NewWorkerRequest("somename")
	//
	//	}
	//}
}

func (p *Parser) workerRequest(name string) {
	defer func() {
		p.log.Printf("[workerRequest][%s] закончил\n", name)
		p.Done()
	}()
	p.log.Printf("[workerRequest][%s] стартанул\n", name)
	var countSleep = 3
	for {
		select {
		case command := <-p.ChanCommand:
			p.log.Printf("[workerRequest][%s] %v\n", command)
			if command == "exit" {
				p.log.Printf("[workerRequest][%s] заканчиваю работу \n", command)
				return
			}
		default:
			p.log.Printf("[workerRequest][%s] длина стока для запросов == `%d`\n", name, len(p.StockSearchRequest))
			if len(p.StockSearchRequest) > 0 {
				sr := p.Pop()
				p.log.Printf("SR: %v\n", sr)
				err := p.getlinkfiles(sr)
				if err != nil {
					p.log.Printf(err.Error())
				}
			} else {
				p.log.Printf("[workerRequest][%s] длина стока для запросов == 0\n", name)
				if countSleep == 3 {
					return
				} else {
					p.log.Printf("[workerRequest][%s] переход в фазу краткого отдыха\n", name)
					time.Sleep(time.Second * 1)
					countSleep++
				}
			}
		}
	}
}

//воркер горутина по работе с результатами поиска
func (p *Parser) workerFile(name string) {
	defer func() {
		p.log.Printf("[workerFile][%s] закончил\n", name)
		p.Done()
	}()
	p.log.Printf("[workerFile][%s] стартанул\n", name)
	//var countSleep = 10
	for {
		select {
		case command := <-p.ChanCommand:
			p.log.Printf("[workerFile][%s] %v\n", command)
			if command == "exit" {
				p.log.Printf("[workerFile][%s] заканчиваю работу \n", command)
				return
			}
		default:
			p.log.Printf("[workerFile][%s] длина стока для запросов == `%d`\n", name, len(p.StockSearchResult))
			if len(p.StockSearchResult) > 0 {
				sr := p.PopResult()
				p.log.Printf("[workerFile %s] SR: %v\n", name, sr)
				err := p.getfile(sr)
				if err != nil {
					p.log.Printf(err.Error())
				}
			} else {
				return
				//p.log.Printf("[workerResult][%s] переход в фазу краткого отдыха\n", name)
				//time.Sleep(time.Second * 3)
				//p.log.Printf("[workerResult][%s] длина стока для запросов == 0\n", name)
				//if countSleep == 10 {
				//	return
				//} else {
				//	p.log.Printf("[workerResult][%s] переход в фазу краткого отдыха\n", name)
				//	time.Sleep(time.Second * 3)
				//	countSleep++
				//}
			}
		}
	}
}

//func (p *Parser) PushStock(url string) {
//	p.StockSearchRequest = append(p.StockSearchRequest, SearchRequest{})
//})
//}
func (p *Parser) Pop() SearchRequest {
	p.Lock()
	ret := (p.StockSearchRequest)[len(p.StockSearchRequest)-1]
	p.StockSearchRequest = (p.StockSearchRequest)[0:len(p.StockSearchRequest)-1]
	p.Unlock()
	return ret
}
func (p *Parser) PopResult() SearchResult {
	p.Lock()
	ret := (p.StockSearchResult)[len(p.StockSearchResult)-1]
	p.StockSearchResult = (p.StockSearchResult)[0:len(p.StockSearchResult)-1]
	p.Unlock()
	return ret
}
func (p *Parser) ShowStockRequestStock() {
	p.requestparser.ShowStock()

}
func (p *Parser) ShowStockRequestURL() {
	p.requestparser.ShowStockRequestURL()

}

//получение файла и сохранение его в базу данных
func (p *Parser) getfile(req SearchResult) (error) {
	//парсю строку запроса для разграничения схем запроса
	ur, err := url.Parse(req.LinkFile)
	if err != nil {
		return err
	}
	resp := &http.Response{}
	if ur.Scheme == "http" {
		resp, err = http.Get(req.LinkFile)
		if err != nil {
			return err
		}
	}
	if ur.Scheme == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		rnew, err := http.NewRequest("GET", req.LinkFile, nil)
		if err != nil {
			fmt.Println(err)
		}
		rnew.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1")
		//делаю запрос для чтения файла
		resp, err = client.Do(rnew)
		if err != nil {
			return err
		}
	}
	//читаю файл
	page, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//запись файл на диск
	namedir := p.transliter.TransliterCyr(req.SReq.RequestText)
	if _, err := os.Stat(namedir); os.IsNotExist(err) {
		os.Mkdir(namedir, os.ModePerm)
	}
	filename := strings.Join([]string{namedir, req.Name}, "/")
	outfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		p.log.Printf(err.Error())
	}
	defer outfile.Close()
	outfile.Write(page)
	return nil

	//запись файла в sqlite3 базу в формате: video BLOB,  img = hash64
	fo := FileObj{
		Name:    req.Name,
		Time:    time.Now().Unix(),
		Ext:     strings.Split(req.Name, ".")[1],
		Size:    0,
		Type:    req.Type,
		Data:    page,
		Request: req.LinkFile,
		Desc:    req.Desc,
	}
	//записываю файл в базу данных
	p.db.InsertRecord(fo)
	p.log.Printf("Успешно записан файл `%s` в базу данных\n", req.Name)
	return nil
}

//обработка запросов к поисковым системам
func (p *Parser) getlinkfiles(req SearchRequest) (error) {
	//создаем https реквестер
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	rnew, err := http.NewRequest("GET", req.RequestURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	rnew.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1")
	b, err := client.Do(rnew)
	if err != nil {
		fmt.Println(err)
	}
	doc, err := goquery.NewDocumentFromReader(b.Body)
	if err != nil {
		p.log.Printf(err.Error())
		return err
	}
	//нахожу все ссылки на файлы в зависимости от типа посиковой системы
	switch req.RequestSearch {
	case SEARCH_TYPE_GOOGLE:
		//p.log.Printf("[getlinkfiles][google] %v\n", doc.Text())
		doc.Find(".rg_meta").Each(func(i int, s *goquery.Selection) {
			g := p.requestparser.google.extractLink(s.Text())
			//создаю новый результат запроса
			nsr := SearchResult{
				SReq:     &req,
				Name:     g.Filename,
				Getting:  false,
				LinkFile: g.SourceFile,
			}
			//добавляю к стоку
			p.Lock()
			p.StockSearchResult = append(p.StockSearchResult, nsr)
			p.Unlock()
		})
	case SEARCH_TYPE_YANDEX:
		p.log.Printf("[getlinkfiles] %v\n", doc.Text())
		doc.Find(".serp-item__link").Each(func(i int, s *goquery.Selection) {
			hr, found := s.Attr("href")
			//fmt.Printf("HREF: `%v`\n", hr)
			if found {
				//извлекаю имя файл + путь
				rp, fn := p.requestparser.yandex.extractLink(hr)
				p.log.Printf("Realpath: [%-100s] Filename: [%-70s] \n", rp, fn)
				//создаю новый результат запроса
				nsr := SearchResult{
					SReq:     &req,
					Name:     fn,
					Getting:  false,
					LinkFile: rp,
				}
				//добавляю к стоку
				p.Lock()
				p.StockSearchResult = append(p.StockSearchResult, nsr)
				p.Unlock()
			} else {
				p.log.Printf("[error found in `href`\n")
			}
		})
	}
	return nil
}
