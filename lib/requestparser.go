//---------------------------------------------------------------------------
//  структура файла запроса
//  тип поисковой глубина;системы;тип запроса по контексту;текст запроса
//	1;google;video;приколы ЖКХ
//  10;yandex;image;приколы ЖКХ
//---------------------------------------------------------------------------

package lib

import (
	"os"
	"io/ioutil"
	"strings"
	"fmt"
	"strconv"
)

type RequestParser struct {
	c      *Config
	p      *Parser
	yandex *yandexRequest
	google *googleRequest
}

func NewRequestParser(c *Config, p *Parser) *RequestParser {
	cc := &RequestParser{
		c:      c,
		p:      p,
		yandex: newYandexRequest(),
		google: newGoogleRequest(),
	}
	cc.parsefile()
	return cc
}
func (r *RequestParser) parsefile() error {
	//открываю файл для разбора'
	fmt.Printf("Requets file : %v\n", r.c.RequestFile)
	f, err := os.Open(r.c.RequestFile)
	if err != nil {
		panic(err)
	}
	//читаю файл
	b, err := ioutil.ReadAll(f)
	//разбиваю по линиям
	lines := strings.Split(string(b), "\n")
	//конвертация
	for _, x := range lines {
		if !strings.HasPrefix("#", x) {
			ar := strings.Split(x, ";")
			if len(ar) == 4 {
				requestSearch := ar[1]
				switch requestSearch {
				case SEARCH_TYPE_GOOGLE:
					sr := &SearchRequest{
						RequestSearch: ar[1],
						RequestType:   ar[2],
						RequestText:   ar[3],
					}
					r.google.makerequest(sr)
					r.p.StockSearchRequest = append(r.p.StockSearchRequest, *sr)
				case SEARCH_TYPE_YANDEX:
					deep, _ := strconv.Atoi(ar[0])
					//в зависимости от глубины создаю нужное количество соответствующих запросов
					for i := 1; i <= deep; i ++ {
						sr := &SearchRequest{
							RequestDeep:   i,
							RequestSearch: ar[1],
							RequestType:   ar[2],
							RequestText:   ar[3],
						}
						r.yandex.makerequest(sr)
						r.p.StockSearchRequest = append(r.p.StockSearchRequest, *sr)
					}
				}
			}
		}
	}
	return nil
}
func (r *RequestParser) makeYandexRequest() {

}
func (r *RequestParser) makeGoogleRequest() {

}
func (r *RequestParser) ShowStock() {
	if len(r.p.StockSearchRequest) == 0 {
		fmt.Printf("Stock is empty\n")
		return
	}
	for _, x := range r.p.StockSearchRequest {
		fmt.Printf("%v:%s:%s:%s\n", x.RequestDeep, x.RequestSearch, x.RequestType, x.RequestText, x.RequestURL)
	}
}
func (r *RequestParser) ShowStockRequestURL() {
	if len(r.p.StockSearchRequest) == 0 {
		fmt.Printf("Stock is empty\n")
		return
	}
	for _, x := range r.p.StockSearchRequest {
		fmt.Printf("[requestparser] %s\n", x.RequestURL)
	}
}
