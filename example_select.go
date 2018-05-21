package main

import (
	"fmt"
	"net/http"
	u "net/url"
	"io/ioutil"
	"sync"
	"time"
	"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/mxk/go-sqlite/sqlite3"
	b64 "encoding/base64"

)

type Stocker struct {
	sync.WaitGroup
	sync.RWMutex
	Stock         []Request
	ResultStock   []Result
	ResultCh      chan Result
	ResultStatus  chan int
	ChanCommand   chan string
	ChanEndWorker chan int
}
type Request struct {
	URL string
}
type Result struct {
	URL    string
	Result []byte
	Error  error
}

func NewStocker() *Stocker {
	return &Stocker{
		ResultCh:      make(chan Result),
		ResultStatus:  make(chan int),
		ChanCommand:   make(chan string),
		ChanEndWorker: make(chan int),
	}
}
func (s *Stocker) Push(url string) {
	s.Stock = append(s.Stock, Request{URL: url})
}
func (s *Stocker) Pop() Request {
	ret := (s.Stock)[len(s.Stock)-1]
	s.Stock = (s.Stock)[0:len(s.Stock)-1]
	return ret
}
func (s *Stocker) LenStock() int {
	return len(s.Stock)
}
func (s *Stocker) Run(countWorker int) {
	defer func() {
		fmt.Printf("== [ the end ] ==\n")
	}()
	var workerEnd = countWorker
	for i := 0; countWorker > 0; countWorker-- {
		go s.worker(fmt.Sprintf("worker_%d", i))
		i++
		s.Add(1)
	}
	fmt.Printf("Awaiting working end...\n")

	for workerEnd > 0 {
		select {
		case result := <-s.ResultCh:
			fmt.Printf("Getting result %v\n", result)
			s.ResultStock = append(s.ResultStock, result)
		case <-s.ChanEndWorker:
			workerEnd--
		default:
			fmt.Printf("Wait result\n")
			time.Sleep(1 * time.Second)
		}
	}
}
func (s *Stocker) worker(name string) {
	if s.LenStock() == 0 {
		fmt.Printf("stock empty\n")
		return
	}
	var counter int = 5
	fmt.Printf(fmt.Sprintf("[%s] Starting \n", name))
	for {
		select {
		case command := <-s.ChanCommand:
			fmt.Printf("[command] %v\n", command)
		default:
			if s.LenStock() > 0 {
				s.Lock()
				url := s.Pop()
				s.Unlock()
				fmt.Printf("Extract %v\n", url)
				result, err := s.Geturl(url)
				if err != nil {
					fmt.Printf(err.Error())
				} else {
					fmt.Printf("Get result success `%s` \n", url.URL)
					s.ResultCh <- result
					fmt.Printf("Send result success `%s` \n", url.URL)
				}
			} else {
				if counter > 0 {
					fmt.Printf(fmt.Sprintf("[%s] Stock is empty, going to sleep\n", name))
					time.Sleep(10 * time.Second)
					counter--
				} else {
					fmt.Printf(fmt.Sprintf("[%s] exnd working\n", name))
					s.Done()
					s.ChanEndWorker <- 1
					return
				}
			}
		}
	}
}
func (s *Stocker) Geturl(req Request) (Result, error) {
	r := Result{
		URL: req.URL,
	}
	client := http.Client{}
	client.Get(req.URL)
	resp, err := client.Get(req.URL)
	if err != nil {
		r.Error = err
		return r, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.Error = err
		return r, err
	}
	r.Result = b
	return r, nil
}

func main() {
	var url = "http://spouk.ru"
	s := NewStocker()
	s.Push(url)
	s.Push("http://eax.ru")
	s.Push("http://python.ru")
	s.Push("http://asm.ru")
	s.Push("http://erlang.ru")
	fmt.Printf("Len LIGO : %v\n", s.LenStock())

	//s.Run(5)

	//---------------------------------------------------------------------------
	//  extract link img
	//---------------------------------------------------------------------------
	var href = "/images/search?p=4&amp;text=%D0%B6%D0%BA%D1%85%20%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B&amp;img_url=http%3A%2F%2Fs02.radikal.ru%2Fi175%2F1707%2F24%2F97a453cd5936.jpg&amp;pos=161&amp;rpt=simage"
	ur, _ := u.Parse(href)
	v, _ := u.ParseQuery(ur.RawQuery)
	fmt.Printf("UR: %v\n", v["img_url"][0])
	link := "http://s02.radikal.ru/i175/1707/24/97a453cd5936.jpg"
	split := strings.Split(link, "/")
	filename := split[len(split)-1:][0]
	fmt.Printf("Filename: %s\n", split[len(split)-1:])

	//get img file from link
	client, _ := http.Get(link)
	page, _ := ioutil.ReadAll(client.Body)
	outfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	outfile.Write(page)
	//writing file to BLOB sqlite3 container
	sql := NewWorkWithSqlite("imgdbs.sqlite3")
	ext := strings.Split(filename, ".")[1]
	f := ImgObject{
		Name:filename,
		Size:100,
		Ext: ext,
		Data:page,
	}
	err = sql.writeSqliteImg(f)
	if err != nil {
		panic(err)
	}

	//---------------------------------------------------------------------------
	//  extract all links from yandex search
	//---------------------------------------------------------------------------
	extracturl := func(link string) (filename string, realpath string) {
		ur, _ := u.Parse(link)
		v, _ := u.ParseQuery(ur.RawQuery)
		realpath = v["img_url"][0]
		split := strings.Split(realpath, "/")
		filename = split[len(split)-1:][0]
		return
	}
	request := "https://yandex.ru/images/search?p=1&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B"
	doc, err := goquery.NewDocument(request)
	if err != nil {
		panic(err)
	}
	doc.Find(".serp-item__link").Each(func(i int, s *goquery.Selection) {
		//hr, _ := s.Find("a").Attr("class")
		hr, found := s.Attr("href")
		if found {
			rp, fn := extracturl(hr)
			fmt.Printf("Realpath: [%-100s] Filename: [%-70s] \n", rp, fn)
		} else {
			fmt.Printf("[error found in `href`\n")
		}
	})

	//---------------------------------------------------------------------------
	//  example work channels
	//---------------------------------------------------------------------------
	//var (
	//	resultChan
	//)
	//for {
	//	select d
	//}

}

type WorkWithSqlite struct {
	dbname string
	db     *sqlite3.Conn
}
type ImgObject struct {
	ID   int64
	Name string
	Ext  string
	Size int64
	Data []byte
}

func NewWorkWithSqlite(dbname string) *WorkWithSqlite {
	sql_table := `
	CREATE TABLE IF NOT EXISTS img(
		Id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		Name TEXT,
		Ext TEXT,
		Size INTEGER,
		Data BLOB,
		InsertedDatetime DATETIME
	);
	`
	//создаю новый инстанс
	d := &WorkWithSqlite{
		dbname: dbname,
	}
	//создаю/открываю базу для манипуляций
	c, err := sqlite3.Open(dbname)
	if err != nil {
		panic(err)
	}
	d.db = c
	//создаю таблицу /если не создана
	err = d.db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
	//возвращаю результат
	return d
}
func (s *WorkWithSqlite) writeSqliteImg(f ImgObject) (error) {
	sql := `
	INSERT OR REPLACE INTO img(
		Name,
		Ext,
		Size,
		Data,
		InsertedDatetime
	) values(?, ?, ?, ?,  CURRENT_TIMESTAMP)
	`
	if s.db == nil {
		panic(fmt.Sprintf("База данных не создана/не открыта\n"))
	}
	sh, err := s.db.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer sh.Close()
	//convert data to base64
	baseConvert := b64.StdEncoding.EncodeToString(f.Data)
	fmt.Printf("Base64Convert: `%v`\n", baseConvert)
	//insert into db
	err = sh.Exec(f.Name, f.Ext, f.Size, baseConvert)
	if err != nil {
		panic(err)
	}
	return nil
}
func (s *WorkWithSqlite) readSqliteImg() (result []ImgObject) {
	sql := `
	SELECT Id, Name, Ext, Size, Data FROM img
	ORDER BY datetime(InsertedDatetime) DESC
	`
	if s.db == nil {
		panic(fmt.Sprintf("База данных не создана/не открыта\n"))
	}
	rows, err := s.db.Query(sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for {
		item := ImgObject{}
		err2 := rows.Scan(&item.ID, &item.Name, &item.Ext, &item.Size, &item.Data)
		if err2 != nil {
			fmt.Printf("Found error in read result from databases \n")
			//panic(err2)
			break
		} else {
			result = append(result, item)
		}
	}
	return
}
