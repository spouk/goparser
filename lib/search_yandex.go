//---------------------------------------------------------------------------
//  враппер для яши для изготовления запросов соответствующего формата
// возращает соответствующий строковый слайс с результатами
// размер слайса равен размеру глубины, глубина в данном случае это количество
// страниц парсинга выдачи
//---------------------------------------------------------------------------

package lib

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

//var (
//	prefix_base    = "https://yandex.ru/images/search=?%s&%s&%s" //основая для запроса
//	prefix_request = "text=%s"                           //текст запроса
//	prefix_deep    = "p=%v"                              //страница запроса
//	prefix_end     = "rpt=%s"                            //type request [image;video]
//)
var (
	prefix_system  = "https://yandex.ru/%s/search?"
	prefix_base    = "https://yandex.ru/%s/search?%s&%s&%s" //основая для запроса
	prefix_request = "text=%s"                              //текст запроса
	prefix_deep    = "p=%v"                                 //страница запроса
	prefix_end     = "rpt=%s"                               //type request [image;video]
)

type yandexRequest struct{}

func newYandexRequest() *yandexRequest {
	return &yandexRequest{}
}
func wrapSearchRequest(s *SearchRequest) {
	nu := new(url.URL)
	u, err := nu.Parse(fmt.Sprintf(prefix_system, s.RequestType))
	if err != nil {
		panic(err)
	}
	params := url.Values{}
	params.Add("text", s.RequestText)
	deep := strconv.Itoa(s.RequestDeep)
	params.Add("p", deep)
	params.Add("rpt", s.RequestType)
	u.RawQuery = params.Encode()
	s.RequestSearch = u.String()
}

func (y *yandexRequest) makerequest(s *SearchRequest) {
	for i := 0; s.RequestDeep > i; i++ {
		nu := new(url.URL)
		u, err := nu.Parse(fmt.Sprintf(prefix_system, s.RequestType))
		if err != nil {
			panic(err)
		}
		params := url.Values{}
		params.Add("text", s.RequestText)
		deep := strconv.Itoa(s.RequestDeep)
		params.Add("p", deep)
		params.Add("rpt", s.RequestType)
		u.RawQuery = params.Encode()
		s.RequestURL = u.String()
	}
}
func (y *yandexRequest) extractLink(link string) (filename string, realpath string){
	ur, _ := url.Parse(link)
	v, err := url.ParseQuery(ur.RawQuery)
	if err != nil {
		fmt.Printf("[yandexRequest][ERROR] %v\n", err.Error())
		return
	} else {
		realpath = v["img_url"][0]
		split := strings.Split(realpath, "/")
		filename = split[len(split)-1:][0]
		fmt.Printf("[REALPATH: `%v` FILENAME: `%v`\n]", realpath, filename)
	}
	return
}