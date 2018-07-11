package library

import (
	"os"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"io"
	"strings"
	"crypto/tls"
	"net/http"
	"math/rand"
	"bufio"
	"log"
	"sync"
	"encoding/json"
	//"time"
	"errors"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
)

type (
	//специфичная структура под google.images объекты
	ResponseStruct struct {
		Isu string `json:"isu"`
		Ity string `json:"ity"`
		Ou  string `json:"ou"`
	}
	Parser struct {
		config       *ConfigTable
		Log          *log.Logger
		stockRequest []*Request
		stockSearch  []*SearchRequest
		sync.WaitGroup
		sync.RWMutex
		chanCommand  chan string
		endChan      chan bool
		//database
		DB *gorm.DB
		//downloader
		Down *Downloader
	}
	//структура описывающая элемент прямой ссылки
	Request struct {
		textRequest string
		linkRequest string
		status      bool
		typeContext string
	}
	//структура описывающая запрос с файла запросов
	SearchRequest struct {
		SearchCore string //google || yandex || etc...
		SearchType string //image || video
		SearchText string //a search text only
		SearchLink string //search request
		//рекурсивный аспект структуру по умолчанию reqursion = false, RecursionIterount = 0
		Reqursion          bool
		RecursionIterCount int //количество итерация, если == 0 то выход
	}
	//структура под базу данных
	DatabaseTable struct {
		ID         int64
		SearchCore string // yandex, google etc
		SearchType string //image || video
		SearchText string //запрос из списка запросов из файла requestlist
		SearchLink string //search request from requestlist
		DirectLink string //прямая ссылка на ресурс
		Active     bool   //скачено или нет
	}
)

//---------------------------------------------------------------------------
//  создание нового инстанса
//---------------------------------------------------------------------------
func NewParser(configFileName string) (*Parser, error) {
	//создаю инстанс парсера
	p := &Parser{
		chanCommand: make(chan string),
		endChan:     make(chan bool),
	}

	//создаю логгер для отладки и вывода сообщений
	p.Log = log.New(os.Stderr, LOGCONFIGPREFIX, LOGCONFIGFLAGS)

	//читаем конфигурационный файл + парсим его на составляющие
	err := p.ConfigReadFromFile(configFileName)
	if err != nil {
		return nil, err
	}

	//читаю файл с запросами с размещением запросов в списке для обработки
	err = p.RequestReadFile()
	if err != nil {
		return nil, err
	}

	//показ содержимого stockSearch
	p.Log.Print(p.stockSearch)

	//создание базы данных
	//create/open database
	db, err := gorm.Open("sqlite3", p.config.Database)
	if err != nil {
		p.Log.Fatal(err)
	}
	//set WAL mode
	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA cache_size = 2000;")
	db.Exec("PRAGMA default_cache_size = 2000;")
	db.Exec("PRAGMA page_size = 4096;")
	db.Exec("PRAGMA synchronous = NORMAL;")
	db.Exec("PRAGMA foreign_keys = on;")
	db.Exec("PRAGMA busy_timeout = 1;")
	//create tables if not  exists
	var listTables = []interface{}{DatabaseTable{}}
	for _, x := range listTables {
		if !db.HasTable(x) {
			err := db.CreateTable(x).Error
			if err != nil {
				log.Print(err)
			}
		}
	}
	//связывем хандлер базы данных
	p.DB = db

	//создаю инстанс загрузчика прямых ссылок видео+картинки
	p.Down = new(Downloader)

	//проверяю путь для сохранения и создаю если его нет
	if _, err := os.Stat(p.config.Pathsave); os.IsNotExist(err) {
		if err = os.Mkdir(p.config.Pathsave, os.ModeDir|os.ModePerm); err != nil {
			p.Log.Fatal(err)
		} else {
			p.Log.Println("Success created directory to save result parse")
		}
	}

	//возвразаю результат
	return p, nil
}

func (p *Parser) Run() {
	p.Log.Printf("STOCK: %v\n", p.stockSearch)
	for _, x := range p.stockSearch {
		fmt.Printf("Stock: %v\n", x)
	}
	//os.Exit(1)

	//p.Add(1)
	//go p.workerDBS()
	//p.Add(1)
	//go p.manager()
	//p.Wait()

	p.Log.Println("STOCKSEARCH: %v\n", p.stockSearch)
	//запуск обработки каждого запроса с паузами обработкой каждого запроса
	for _, e := range p.stockSearch {
		switch e.SearchType {
		case "video":
			err := p.GoogleGetLinksVideo(e)
			if err != nil {
				fmt.Printf("[video] Error: %v\n", err.Error())
				//time.Sleep(time.Minute * 60)
			} else {
				//time.Sleep(time.Minute * 30)
			}
		case "image":
			err := p.GoogleGetLinksImage(e)
			if err != nil {
				fmt.Printf("[image] Error: %v\n", err.Error())
				//time.Sleep(time.Minute * 60)
			} else {
				//time.Sleep(time.Minute * 30)
			}
		default:
			p.Log.Printf("ошибка в типе запроса контекста к поисковой системе, должно быть `video` или `image`")
		}
	}
	//вывод результата
	p.showRequestStock()

	//запуск горутин для обработки прямых ссылок
	for i, x := range p.stockRequest {
		switch x.typeContext {
		case "video":
			p.Add(1)
			go p.DownloaderVideo(fmt.Sprintf("[VIDEO#%d]", i), p.config.Pathsave, x.linkRequest, true)
		case "image":
			p.Add(1)
			go p.DownloadImage(p.config.Pathsave, x.linkRequest, true)
		default:
			log.Printf("WRONG TYPE REQUEST `%v`:%v\n", x.typeContext, x)
		}
	}
	//ожидание завершения работы горутин
	p.Wait()
	p.Log.Printf("All SUCCESS DOWNLOAD LINKS")


	p.Log.Print("Parser done working\n")
	return
}

//---------------------------------------------------------------------------
//  MANAGER
//---------------------------------------------------------------------------
func (p *Parser) manager() {
	defer func() {
		p.Done()
		close(p.endChan)
	}()
	p.Log.Println("STOCKSEARCH: %v\n", p.stockSearch)
	//запуск обработки каждого запроса с паузами обработкой каждого запроса
	for _, e := range p.stockSearch {
		switch e.SearchType {
		case "video":
			err := p.GoogleGetLinksVideo(e)
			if err != nil {
				fmt.Printf("[video] Error: %v\n", err.Error())
				//time.Sleep(time.Minute * 60)
			} else {
				//time.Sleep(time.Minute * 30)
			}
		case "image":
			err := p.GoogleGetLinksImage(e)
			if err != nil {
				fmt.Printf("[image] Error: %v\n", err.Error())
				//time.Sleep(time.Minute * 60)
			} else {
				//time.Sleep(time.Minute * 30)
			}
		default:
			p.Log.Printf("ошибка в типе запроса контекста к поисковой системе, должно быть `video` или `image`")
		}
	}
	//вывод результата
	p.showRequestStock()

	//запуск горутин для обработки прямых ссылок
	for i, x := range p.stockRequest {
		switch x.typeContext {
		case "video":
			p.Add(1)
			go p.Down.DownloaderVideo(fmt.Sprintf("[VIDEO#%d]", i), p.config.Pathsave, x.linkRequest, true, p.WaitGroup)
		case "image":
			p.Add(1)
			go p.Down.DownloadImage(p.config.Pathsave, x.linkRequest, true, p.WaitGroup)
		default:
			log.Printf("WRONG TYPE REQUEST `%v`:%v\n", x.typeContext, x)
		}
	}
	//ожидание завершения работы горутин
	p.Wait()
	p.Log.Printf("All SUCCESS DOWNLOAD LINKS")
	return
}
func (p *Parser) worker(id int) {
	p.Log.Printf("worker#%d starting...\n", id)
	defer func() {
		p.Done()
		p.Log.Printf("worker#%d end work\n", id)
	}()
	for {
		select {
		case command := <-p.chanCommand:
			if command == "exit" {
				return
			}
		case <-p.endChan:
			return
		default:
			if len(p.stockRequest) > 0 {
				p.Lock()
				element := p.stockSearch[0]
				p.stockSearch = append(p.stockSearch[:0], p.stockSearch[1:]...)
				p.Unlock()
				switch element.SearchType {
				case "video":
				case "image":
				default:
					p.Log.Printf("Wrong searchtype `%v\n`", element.SearchType)
				}

			} else {
				return
			}
		}
	}
}
func (p *Parser) workerDBS() {
	defer func() {
		p.Done()
	}()
	for {
		select {
		case <-p.endChan:
			return
		default:
			if len(p.stockRequest) > 0 {
				var element = p.stockRequest[0]
				p.Lock()
				p.stockRequest = append(p.stockRequest[:0], p.stockRequest[1:]...)
				p.Unlock()
				var records []DatabaseTable
				if err := p.DB.Find(&records).Error; err != nil {
					p.Log.Println(err)
				} else {
					if func(x *Request, records []DatabaseTable) bool {
						for _, z := range records {
							if x.linkRequest == z.DirectLink {
								return true
							}
						}
						return false
					}(element, records) == false {
						var newRecord = &DatabaseTable{
							DirectLink: element.linkRequest,
							Active:     false,
							SearchLink: element.textRequest,
							SearchType: element.typeContext,
							SearchText: element.textRequest,
						}
						if err := p.DB.Create(newRecord).Error; err != nil {
							p.Log.Println(err)
						}
					} else {
						p.Log.Printf("Found Dublicate DIRECT LINK IN DATABASE `%v`\n", element.linkRequest)
					}
				}
			}
		}
	}
}

//---------------------------------------------------------------------------
//  GOOGLE:VIDEO
//---------------------------------------------------------------------------
//список страниц выдачи запроса - получение ссылок
func (p *Parser) GoogleGetLinksVideo(r *SearchRequest) (error) {
	p.Log.Printf("starting video reqvester\n")

	//получаю список доступных страниц на выдаче поисковой системы
	resp, err := p.MakeRequestSearchSystem(r.SearchLink)
	if err != nil {
		return err
	}
	//проверка на http код
	if resp.StatusCode != 200 {
		//ошибка или отлуп идли блокировка
		fmt.Printf("CODE: %v : %v\n", resp.StatusCode, resp.Status)
		return errors.New("Error: HTTP Service 403/503 errors...")
	}

	//создаю инстанс ридера для корректной обработки в поисковике
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	//нахожу все ссылки на имеющиеся страницы по первому запросу, вида /search++++
	//для создания корректной ссылки требуется название поисковой системы + часть поискового запроса
	var stock []string
	doc.Find("table#nav a.fl").Each(func(i int, l *goquery.Selection) {
		str, exists := l.Attr("href")
		p.Log.Println("[href] LINK VIDEO: %v\n", GOOGLEBASA+str)
		if exists {
			stock = append(stock, GOOGLEBASA+str)
			p.Log.Println("LINK VIDEO: %v\n", GOOGLEBASA+str)
		}
	})
	fmt.Printf("[%d] STOCKVIDEO: %v\n", len(stock), stock)

	//обработка списка ссылок страниц с выдачей
	for _, x := range stock {
		p.parseGoogleVideoLinks(x, r)
	}

	return nil
}

//парсит все ссылки на странице выдачи с гугла по видео запросу
func (p *Parser) parseGoogleVideoLinks(req string, s *SearchRequest) (error) {
	//при возврате возвращаем триггер на горутину для WaitGroup
	//defer func() {
	//	p.Done()
	//}()
	//получаю список доступных страниц на выдаче поисковой системы
	resp, err := p.MakeRequestSearchSystem(req)
	if err != nil {
		return err
	}
	//check http error
	if resp.StatusCode != 200 {
		return errors.New("Error: HTTP Service 403/503 errors...")
	}
	//создаю инстанс ридера для корректной обработки в поисковике
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	//ищу все блоки ссылками на вdiaидосы
	doc.Find("h3.r a").Each(func(i int, l *goquery.Selection) {
		str, exists := l.Attr("href")
		if exists {
			//парсю название ресурса для фильтрации результата
			u, err := url.Parse(str)
			if err != nil {
				fmt.Printf("error: %v"+
					""+
					"\n", err.Error())
			} else {
				//наш хост - `труба`
				if u.Host == "www.youtube.com" {
					fmt.Printf("[%s] Result link video: %v\n", u.Host, str)
					p.Lock()
					p.stockRequest = append(p.stockRequest, &Request{linkRequest: str, textRequest: s.SearchText, status: false, typeContext: s.SearchType})
					p.Unlock()
				}
			}
		}
	})
	//debug msg
	p.Log.Printf("back from `parseGoogleVideo`\n")

	//возвращаю результат
	return nil
}

//---------------------------------------------------------------------------
//  GOOGLE:IMAGE
//---------------------------------------------------------------------------
//получение ссылок по запросу по картинкам
func (p *Parser) GoogleGetLinksImage(r *SearchRequest) (error) {
	p.Log.Printf("Image parser gorouti: %v\n", r)
	//получаю список доступных страниц на выдаче поисковой системы
	resp, err := p.MakeRequestSearchSystem(r.SearchLink)
	if err != nil {
		return err
	}
	//создаю инстанс ридера для корректной обработки в поисковике
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	//ищу все диваны с json вставкой
	doc.Find("div.rg_meta").Each(func(i int, l *goquery.Selection) {
		rs := new(ResponseStruct)
		err := json.Unmarshal([]byte(l.Text()), rs)
		if err != nil {
			fmt.Printf("Error ---> `%v` %v\n", err.Error(), rs)
		} else {
			p.Lock()
			p.stockRequest = append(p.stockRequest, &Request{linkRequest: rs.Ou, textRequest: r.SearchText, status: false, typeContext: r.SearchType})
			p.Unlock()
		}
	})

	return nil
}

//читаю конфигурационный файл и возвращаю инстанс структуру
func (p *Parser) ConfigReadFromFile(configFilename string) (error) {
	if err, ci := p.newConfigRead(configFilename); err != nil {
		return err
	} else {
		p.config = ci
		return nil
	}
}

//читаю файл с поисковыми запросами
func (p *Parser) RequestReadFile() (error) {
	var (
		str      []byte
		stockStr []string
		result   []*SearchRequest
	)
	if file, err := os.Open(p.config.RequestFile); err != nil {
		return err
	} else {
		defer file.Close()
		b := bufio.NewReader(file)
		//читаю файл посредством чтения каждой строки с размещением полученного в списке строк
		for io.EOF != err {
			str, _, err = b.ReadLine()
			stockStr = append(stockStr, string(str))
		}
		//SearchCore string //google || yandex || etc...
		//SearchType string //image || video
		//SearchText string //a search text only
		//SearchLink string //search request

		//сортирую полученные строки
		for _, l := range stockStr {
			if len(l) > 0 {
				if strings.HasPrefix(l, ":::") {
					splitter := strings.Split(l, ":::")
					if len(splitter) == 5 {
						rr := &SearchRequest{
							SearchCore: splitter[1],
							SearchType: splitter[2],
							SearchText: splitter[3],
							SearchLink: splitter[4]}
						result = append(result, rr)
					}
				}
			}

		}
		p.stockSearch = result
		return nil
	}
}

//функция генерируют рандомный заголовок для запроса к поисковой системе
func (p *Parser) randomMakeHeader() string {
	return listHeaders[rand.Intn(len(listHeaders))]
}

//функция делает запрос к поисковой системе с целью получения ответа для передачи далее на обработку
func (p *Parser) MakeRequestSearchSystem(r string) (*http.Response, error) {
	//формирую новый http запрос
	req, err := http.NewRequest("GET", r, nil)
	//если ошибка возврат из функции
	if err != nil {
		return nil, err
	}
	//создаю транспорт с поддержкой TLS для корректной обработки подключений по протоколу HTTPS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//создаю http клиента
	client := &http.Client{Transport: tr}
	//добавляю корректный заголовок для создания иллюзии запроса от браузера
	req.Header.Set("User-Agent", p.randomMakeHeader())
	//осуществляю запрос
	b, err := client.Do(req)
	//если возникает ошибка  при осуществлении запроса к сайту, то получаемая ошибка возвращается назад
	//после выхода из функции
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (p *Parser) showRequestStock() {
	if len(p.stockRequest) == 0 {
		p.Log.Printf("StockRequest is empty\n")
		return
	}
	for i, x := range p.stockRequest {
		p.Log.Printf("[%d] %s\n", i, x.linkRequest)
	}
}
