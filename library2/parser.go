package library2

import (
	"sync"
	"io"
	"log"
)

const (
	YANDEXBASE      = "https://yandex.ru"
	YANDEXBASEIMAGE = "https://yandex.ru/images/search?p=%d&text=%s"
	YANDEXBASEVIDEO = "https://yandex.ru/video/search?p=%d&text=%s"

	//YANDEXBASEIMAGE = "https://yandex.ru/images/search?p=1&text=&rpt=image"
	//"https://yandex.ru/images/search?p=1&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B"
	//https://yandex.ru/video/search?p=2&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
	//https://yandex.ru/video/search?p=3&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
	//https://yandex.ru/video/search?p=4&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
	//https://yandex.ru/video/search?p=4&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
)
const (
	PREFIX = "[parser] "
)

type Parser struct {
	sync.WaitGroup
	Stock []Box
	Logger *log.Logger
}
type Box struct {
	Request string
	Count   int
	Value   io.Writer
	Name    string
}
func New(logout io.Writer) *Parser {
	return &Parser{
		Logger: log.New(logout, PREFIX, log.Ldate  | log.Lshortfile),
	}
}
func (p *Parser) worker() {
	if len(p.Stock)  == 0 {
		p.Logger.Printf("Стек пустой, нечего обрабатывать\n")
		return
	}
	for true {

	}
}
