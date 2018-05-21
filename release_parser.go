//---------------------------------------------------------------------------
//  yandex search parametrs
// lr=номер региона поиска, к примеру Тула = 15, без учета региона поиска = 0 (lr=0)
// p=номер страницы выдачи
// text=текст запроса
//---------------------------------------------------------------------------

package main

import (
	"sync"
	"log"
	"time"
	"os"
	"io/ioutil"
	"github.com/go-yaml/yaml"
	"bufio"
	"strings"
	"strconv"
	"fmt"
	"io"
	"crypto/tls"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

type (
	configRelease struct {
		DB          string        `yaml:"db"`
		LogfileName string        `yaml:"logfilename"`
		WorkRequest int           `yaml:"workrequest"`
		WorkPause   time.Duration `yaml:"workpause"`
		FileRequest string        `yaml:"filerequest"`
		RequestFile string        `yaml:"requestfile"`
	}
	request struct {
		count   int
		types   string
		context string
		request string
	}
	Parser struct {
		sync.RWMutex
		sync.WaitGroup
		log         *log.Logger
		config      *configRelease
		StockRequst []request
	}
)

const (
	LOGCONFIGPREFIX = "[parser]"
	LOGCONFIGFLAGS  = log.Lshortfile | log.Ldate | log.Ltime
	yandexRequest   = "https://yandex.ru/search"
)

var p *Parser

func init() {
	var (
		configName = "config.yml"
	)

	//make instance
	p = NewParserRelease()
	//read config
	if err := p.readConfig(configName); err != nil {
		p.log.Panic(err)
	}
	//make log file
	if p.config.LogfileName != "" {
		f, err := os.OpenFile(p.config.LogfileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			p.log = log.New(os.Stdout, LOGCONFIGPREFIX, LOGCONFIGFLAGS)
			p.log.Printf(err.Error())
		} else {
			p.log = log.New(f, LOGCONFIGPREFIX, LOGCONFIGFLAGS)
		}
	} else {
		p.log = log.New(os.Stdout, LOGCONFIGPREFIX, LOGCONFIGFLAGS)
	}
	//read requestfile
	err := p.readRequestFile(p.config.RequestFile)
	if err != nil {
		p.log.Panic(err)
	}
	fmt.Println("all ok init")
}
func main() {
	fmt.Printf("Config: %v: %s\n", p.config, p.config.RequestFile)
	fmt.Printf("StockRequest: %v\n", p.StockRequst)
	for _, x := range p.StockRequst {
		fmt.Println(x)
	}
	p.parserYandeLinks(p.StockRequst[2])
}
func NewParserRelease() *Parser {
	return new(Parser)
}
func (p *Parser) readConfig(configFileName string) (error) {
	//открываю файл с конфигом для чтения
	f, err := os.Open(configFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	//парсю файл
	p.config = new(configRelease)
	err = yaml.Unmarshal(b, p.config)
	if err != nil {
		return err
	}
	return nil
}
func (p *Parser) readRequestFile(filename string) (error) {

	//set variables
	var (
		line      []byte
		err       error = nil
		stockLine []string
	)
	//try open file
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("err open file: %v\n", err)
		return err
	}
	//close after return function
	defer f.Close()
	//make reader
	b := bufio.NewReader(f)
	//read lines
	for io.EOF != err {
		line, _, err = b.ReadLine()
		stockLine = append(stockLine, string(line))
	}
	fmt.Printf("StockLine: %v\n", stockLine)
	//check `true` line request and trash ex. comments etc...
	for _, x := range stockLine {
		arr := strings.Split(x, ";")
		if arr[0] == "REQ" {
			c, err := strconv.Atoi(arr[1])
			if err != nil {
				p.log.Printf(err.Error())
			} else {
				p.StockRequst = append(p.StockRequst, request{count: c, types: arr[2], context: arr[3], request: arr[4]})
			}
		}
	}
	return nil
}
func (p *Parser) parserYandeLinks(r request) (error) {
	//make valid http string for correct request
	req, err := http.NewRequest("GET", yandexRequest, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("lr", "0")
	q.Add("text", r.request)
	req.URL.RawQuery = q.Encode()

	//make http instance + set TLS for correction https connections
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//make http client
	client := &http.Client{Transport: tr}

	//make header for `stealth` request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1")

	//make request
	b, err := client.Do(req)
	if err != nil {
		return err
	}
	//read answer from resource
	//buf, _ := ioutil.ReadAll(b.Body)

	//find all links in answer
	doc, err := goquery.NewDocumentFromReader(b.Body)
	if err != nil {
		return err
	}
	doc.Find("a").Each(func(i int, l *goquery.Selection) {
		hr, _:= l.Attr("href")
		fmt.Printf("found LINK %v\n", hr)
	})

	//find all links
	//doc.Find(".link_outer_yes").Each(func(i int, s *goquery.Selection) {
	doc.Find(".serp-item").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Found element %v\n", s)
		//rr , _ := s.Attr("href")
		//fmt.Printf("Link: %s\n", rr)

		//g := &GoogleImage{}
		//err := json.Unmarshal([]byte(s.Text()), g)
		//if err != nil {
		//	fmt.Printf("Error unmarshal: %v\n", err.Error())
		//} else {
		//	fmt.Printf("link file: %s\n", g.SourceFile)
		//}
		s.Find("a").Each(func(i int, l *goquery.Selection) {
			hr, _:= l.Attr("href")
			fmt.Printf("found LINK %v\n", hr)
		})
	})
	return nil
}
