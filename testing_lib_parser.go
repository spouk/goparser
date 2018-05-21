package main

import (
	"yuristlic.parser/lib"
	"os"
	"fmt"
	"crypto/tls"
	"net/http"
	"github.com/PuerkitoBio/goquery"

	"encoding/json"
	"io/ioutil"
	"bytes"
	"time"
	"sync"
)

var (
	config      = "config.yml"
	requestfile = "requestlist.lst"
)

type ExampleStock struct {
	sync.RWMutex
	sync.WaitGroup
	timer *time.Ticker
	Stock []int
}

func NewExampleStock(timeSecond int) *ExampleStock {
	return &ExampleStock{
		timer: time.NewTicker(time.Second * time.Duration(timeSecond)),
	}
}
func (s *ExampleStock) pop() int {
	s.Lock()
	ret := (s.Stock)[len(s.Stock)-1]
	s.Stock = (s.Stock)[0:len(s.Stock)-1]
	s.Unlock()
	return ret
}
func (s *ExampleStock) workerTimer(name string) {
	defer func() {
		s.Done()
	}()
	fmt.Printf("workertimer start %s\n", name)
	//	var count = 100
	for {
		select {
		case _ = <-s.timer.C:
			fmt.Printf("[%s] get new timer: `%d`\n", name, s.pop())
		default:
			if len(s.Stock) == 0 {
				return
			}
			//if count == 0 {
			//	fmt.Printf("[%s] go exit, cya\n", name)
			//	return
			//} else {
			//	//fmt.Printf("[%s] going sleep\n", name)
			//	//time.Sleep(time.Millisecond * 500)
			//	count --
			//}
		}
	}
}
func (s *ExampleStock) Run(countWorkers int) {
	for i := 0; i <= countWorkers; i ++ {
		go s.workerTimer(fmt.Sprintf("WORKER#%d", i))
	}
	s.Add(countWorkers)
	s.Wait()
}
func main() {
	//---------------------------------------------------------------------------
	//  эксперименты с таймерами и каналами
	//---------------------------------------------------------------------------
	//создаю новый стек с данными
	stock := NewExampleStock(2)
	for i := 0; i < 300; i++ {
		stock.Stock = append(stock.Stock, i)
	}

	//stock.Run(10)

	//return

	//создаем парсер
	p := lib.NewParser(config, os.Stdout)
	fmt.Printf("ParserExample: %v\n", p)
	fmt.Printf("Stock parser: %v\n", p.StockSearchRequest)
	//p.ShowStockRequestStock()
	//p.ShowStockRequestURL()
	for _, x := range p.StockSearchRequest {
		fmt.Printf("==> %v\n", x.RequestURL)
	}
	p.Run(2, 10)
	return

	//---------------------------------------------------------------------------
	//  testing google
	//---------------------------------------------------------------------------
	//request := "https://www.google.ru/search?q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%96%D0%9A%D0%A5%0D&tbm=isch"
	//request2 := "https://www.google.ru/search?safe=off&biw=1096&bih=548&tbm=isch&sa=1&q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%B2+%D0%B6%D0%BA%D1%85&oq=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%B2+%D0%B6%D0%BA%D1%85&gs_l=psy-ab.12...0.0.0.871550.0.0.0.0.0.0.0.0..0.0....0...1..64.psy-ab..0.0.0....0.A4p0MsydTuM"
	//requestgoogle(request2)
	//return

	f, err := os.Open("google_answer_noformated.html")
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	requestfromByteArray(buf)

}

type GoogleImage struct {
	SourceSite    string `json:"isu"`
	Extfile       string `json:"ity"`
	SourceFile    string `json:"ou"`
	Size          int    `json:"ow"`
	Desc          string `json:"s"`
	EncryptGoogle string `json:"tu"`
}

func requestfromByteArray(b []byte) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	result := []*GoogleImage{}
	doc.Find(".rg_meta").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Found element %v\n", s.Text())
		g := &GoogleImage{}
		err := json.Unmarshal([]byte(s.Text()), g)
		if err != nil {
			fmt.Printf("Error unmarshal: %v\n", err.Error())
		} else {
			fmt.Printf("link file: %s\nDesc: `%s`\n", g.SourceFile, g.Desc)
			result = append(result, g)
		}
		//s.Find("a").Each(func(i int, l *goquery.Selection) {
		//	hr, _:= l.Attr("href")
		//	fmt.Printf("found LINK %v\n", hr)
		//})
	})
	fmt.Printf("Total: %d\n", len(result))

}
func requestgoogle(r string) {

	//создаем https реквестер
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.1")
	b, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	buf, _ := ioutil.ReadAll(b.Body)
	fmt.Printf("Result reqest: %v\n", string(buf))
	return
	//нахожу все элементы в ответе
	doc, err := goquery.NewDocumentFromReader(b.Body)
	if err != nil {
		panic(err)
	}
	//нахожу все ссылки на файлы в зависимости от типа посиковой системы
	doc.Find(".rg_meta").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Found element %v\n", s)
		g := &GoogleImage{}
		err := json.Unmarshal([]byte(s.Text()), g)
		if err != nil {
			fmt.Printf("Error unmarshal: %v\n", err.Error())
		} else {
			fmt.Printf("link file: %s\n", g.SourceFile)
		}
		//s.Find("a").Each(func(i int, l *goquery.Selection) {
		//	hr, _:= l.Attr("href")
		//	fmt.Printf("found LINK %v\n", hr)
		//})
	})
}

//fmt.Printf("FOUND ITEM: %v: %v \n", s, s.Text())
//hr, found := s.Attr("href")
//fmt.Printf("HREF: `%v`\n", hr)
//if found {
//	//извлекаю имя файл + путь
//	rp, fn := p.requestparser.google.extractLink(hr)
//	p.log.Printf("[google] Realpath: [%-100s] Filename: [%-70s] \n", rp, fn)
//	//создаю новый результат запроса
//	nsr := SearchResult{
//		Name:     fn,
//		Getting:  false,
//		LinkFile: rp,
//	}
//	//добавляю к стоку
//	p.Lock()
//	p.StockSearchResult = append(p.StockSearchResult, nsr)
//	p.Unlock()
//} else {
//	p.log.Printf("[error found in `href`\n")
//}

//buf, _ :=ioutil.ReadAll(b.Body)
//fmt.Printf("Result reqest: %v\n", string(buf))
//---------------------------------------------------------------------------
//  тестирую яшу и гоуртины с интервалами между запросами
//---------------------------------------------------------------------------
func workeryandex(name string) {
	fmt.Printf("starting worker#`%s`\n", name)
	for {
		select {
		default:
			//проверка на длину стека
			//проверка на таймер, если истек
			//
		}
	}
}
func yandexrequest(countWoker int) {

}
