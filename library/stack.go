package library

import (
	"sync"
	"fmt"
	"time"
	"log"
	"os"
	"net/http"
	"golang.org/x/net/html"
	"crypto/tls"
)

const (
	STACKPREFIX = "[parser]"
)

type (
	Stack struct {
		//base stack
		sync.Mutex
		sync.WaitGroup
		Stock            map[StockKey][]StockResult
		StatusCh         chan int
		TypeRequest      TypeRequest
		TextRequest      string
		TypeSearchEngine TypeSearchEngine
		DepthRequest     int
		//header
		Headers Headers
		//logger
		logger *log.Logger
	}
	StockKey struct {
		Date         time.Time
		TextRequest  string
		TypeRequest  TypeRequest
		SearchEngine TypeSearchEngine
		Depth        int
	}
	StockResult struct {
		Page  int
		Key   string
		Value string
	}
	Worker struct {
		Name     string
		StatusCh chan int
		Stock    *map[StockKey][]StockResult
		SW       *sync.WaitGroup
	}
	TypeRequest struct {
		Image bool
		Video bool
	}
	TypeSearchEngine struct {
		Google bool
		Yandex bool
	}
)

func NewWorker(name string, statusCh chan int, stock *map[StockKey][]StockResult, sw *sync.WaitGroup) *Worker {
	w := new(Worker)
	w.Name = name
	w.StatusCh = statusCh
	w.Stock = stock
	w.SW = sw
	return w
}
func (w *Worker) Run() {
	defer func() {
		w.SW.Done()
	}()

	time.Sleep(1000 * time.Millisecond)
	fmt.Printf("%s the end \n", w.Name)
	return
}

func NewStack(textRequest string, typeRequest TypeRequest, typeSearchEngine TypeSearchEngine, depth int) *Stack {
	//make Stack
	s := &Stack{}
	s.Stock = make(map[StockKey][]StockResult)
	s.TextRequest = textRequest
	s.TypeRequest = typeRequest
	s.TypeSearchEngine = typeSearchEngine
	s.DepthRequest = depth
	s.logger = log.New(os.Stdout, STACKPREFIX, log.Ltime|log.Ldate|log.Lshortfile)
	return s
}
func (s *Stack) Run(countWorker int) {
	//make request pool
	//for x := 0; x <= s.DepthRequest; x++ {
	//	s.Stock[]
	//}
	defer func() {
		fmt.Printf("The end all workers\n")
	}()
	for x := 1; x < countWorker; x ++ {
		s.Add(1)
		fmt.Printf(fmt.Sprintf("Run Worker %d\n", x))
		nw := NewWorker(
			fmt.Sprintf("Worker %d", x),
			s.StatusCh,
			&s.Stock,
			&s.WaitGroup,
		)
		go nw.Run()
	}
	s.Wait()
	return
}

//---------------------------------------------------------------------------
//  function for random generate result for testing
//---------------------------------------------------------------------------
func (s *Stack) RandomGenerate(count int, typeRequest TypeRequest, textReqeust string, SerchEngine TypeSearchEngine) {
	//make new key
	nk := StockKey{
		TypeRequest:  typeRequest,
		TextRequest:  textReqeust,
		SearchEngine: SerchEngine,
		Date:         time.Now(),
		Depth:        count,
	}
	//generate random results
	for x := 0; x < count; x++ {
		nr := StockResult{
			Key:  fmt.Sprintf("KEY %d", x),
			Page: 1,
		}
		s.Stock[nk] = append(s.Stock[nk], nr)
	}
}
func (s *Stack) ShowStock() {
	for k, v := range s.Stock {
		fmt.Printf("[%v] [%v] \n", k, v)
	}
}

//---------------------------------------------------------------------------
//  HTTP PARSER
//---------------------------------------------------------------------------
func (s *Stack) HttpRequest(url string) (bool, error) {
	//делаю секретный транспорт для запросов по https
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//делаю новый запрос
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Printf(err.Error())
		return false, err
	}
	//добавляю к запросу рандомный заголовок
	s.Headers.AssignHeader(request)
	//создаю клиента для передачи запроса
	client := &http.Client{Transport: tr}
	//осуществляю запрос
	resp, err := client.Do(request)
	if err != nil {
		s.logger.Printf(err.Error())
		return false, err
	}
	////читаю тело ответа
	//b , err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	s.logger.Printf(err.Error())
	//	return false, err
	//}
	//парсю тело ответа в поиске ссылок
	t := html.NewTokenizer(resp.Body)
	for {
		tt := t.Next()
		switch {
		case tt == html.ErrorToken:
			return false, nil
		case tt == html.StartTagToken:
			t1 := t.Token()
			linkA := t1.Data == "a"
			if linkA {
				s.logger.Printf("Found link: %v\nAttrs: %v\n", t1.Data, t1.Attr)

				for _, a := range t1.Attr {
					if a.Key == "href" {
						s.logger.Printf("[%v] Href: %v\n", t1.Data, a.Val)
						break
					}

				}
			}
		}
	}
}
