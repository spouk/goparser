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
	"time"
	"errors"
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
)

//---------------------------------------------------------------------------
//  создание нового инстанса
//---------------------------------------------------------------------------
func NewParser(configFileName string) (*Parser, error) {
	//создаю инстанс парсера
	p := new(Parser)

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

	return p, nil
}

func (p *Parser) Run() {
	p.Add(1)
	go p.manager()
	p.Wait()
	p.Log.Print("Parser done working\n")
	return
}

//---------------------------------------------------------------------------
//  MANAGER
//---------------------------------------------------------------------------
func (p *Parser) manager() {
	defer func() {
		p.Done()
	}()
	//запуск обработки каждого запроса с паузами обработкой каждого запроса
	for _, e := range p.stockSearch {
		switch e.SearchType {
		case "video":
			err := p.GoogleGetLinksVideo(e)
			if err != nil {
				fmt.Printf("[video] Error: %v\n", err.Error())
				time.Sleep(time.Minute * 60)
			} else {
				time.Sleep(time.Minute * 30)
			}
		case "image":
			err := p.GoogleGetLinksImage(e)
			if err != nil {
				fmt.Printf("[image] Error: %v\n", err.Error())
				time.Sleep(time.Minute * 60)
			} else {
				time.Sleep(time.Minute * 30)
			}
		default:
			p.Log.Printf("ошибка в типе запроса контекста к поисковой системе, должно быть `video` или `image`")
		}
	}
	//вывод результата
	p.showRequestStock()

	return
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
	doc.Find("table#nav a").Each(func(i int, l *goquery.Selection) {
		str, exists := l.Attr("href")
		if exists {
			stock = append(stock, GOOGLEBASA+str)
		}
	})

	//обработка списка ссылок страниц с выдачей
	for _, x := range stock {
		p.parseGoogleVideoLinks(x, r)
	}

	return nil
}

//парсит все ссылки на странице выдачи с гугла по видео запросу
func (p *Parser) parseGoogleVideoLinks(req string, s *SearchRequest) (error) {
	//при возврате возвращаем триггер на горутину для WaitGroup
	defer func() {
		p.Done()
	}()
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
