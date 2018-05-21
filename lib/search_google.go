//---------------------------------------------------------------------------
//  враппер для гугла для изготовления запросов соответствующего формата
// возращает соответствующий строковый слайс с результатами
// размер слайса равен размеру глубины, глубина в данном случае это количество
// страниц парсинга выдачи
//---------------------------------------------------------------------------

package lib

import (
	"fmt"
	"net/url"

	"encoding/json"
	"strings"
)

//https://www.google.ru/search?q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%B2+%D0%B6%D0%BA%D1%85&num=100&safe=off&source=lnms&tbm=isch&sa=X&ved=0ahUKEwizgPPWivXWAhXsIJoKHR4oBmcQ_AUICigB&biw=1096&bih=548
//https://www.google.ru/search?q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%B2+%D0%B6%D0%BA%D1%85&tbm=isch&num=100
var (
	prefix_system_google  = "https://www.google.ru/search?"
	prefix_base_google    = "https://www.google.ru/search?%s&%s&num=1000" //основая для запроса
	prefix_request_google = "q=%s"                                        //текст запроса
	//prefix_deep_google    = "p=%v"                                 //страница запроса //в гугле нет
	prefix_end_google = "tbm=%s" //type request [image;video] [images == `tbm=isch`] [video='tbm=vid']

)

type googleRequest struct{}
type GoogleImage struct {
	SourceSite    string `json:"isu"`
	Extfile       string `json:"ity"`
	SourceFile    string `json:"ou"`
	Size          int    `json:"ow"`
	Desc          string `json:"s"`
	EncryptGoogle string `json:"tu"`
	Filename       string
}

func newGoogleRequest() *googleRequest {
	return &googleRequest{}
}

func (y *googleRequest) makerequest(s *SearchRequest) {
	nu := new(url.URL)
	u, err := nu.Parse(prefix_system_google)
	if err != nil {
		panic(err)
	}
	params := url.Values{}
	params.Add("q", s.RequestText) //текст запроса
	switch s.RequestType {
	case SEARCH_TYPE_CONTEXT_IMAGES:
		s.RequestType = "isch"
	case SEARCH_TYPE_CONTEXT_VIDEO:
		s.RequestType = "vid"
	}
	params.Add("tbm", s.RequestType) //тип запроса
	u.RawQuery = params.Encode()
	s.RequestURL = u.String()
	fmt.Printf("GOOGLE RERQUEST: %v\n", s)
}
func (y *googleRequest) extractLink(str string) (g *GoogleImage) {
	g = &GoogleImage{}
	err := json.Unmarshal([]byte(str), g)
	if err != nil {
		fmt.Printf("[googleextract] Error unmarshal: %v\n", err.Error())
	} else {
		fmt.Printf("link file: %s\nDesc: `%s`\n", g.SourceFile, g.Desc)
		ss := strings.Split(g.SourceFile, "/")
		g.Filename = ss[len(ss) - 1]
		return
	}
	//ur, _ := url.Parse(link)
	//v, err := url.ParseQuery(ur.RawQuery)
	//if err != nil {
	//	fmt.Printf("[googleRequest][ERROR] %v\n", err.Error())
	//	return
	//} else {
	//	realpath = v["imgurl"][0]
	//	split := strings.Split(realpath, "/")
	//	filename = split[len(split)-1:][0]
	//	fmt.Printf("[googleRequest] [REALPATH: `%v` FILENAME: `%v`\n]", realpath, filename)
	//}
	return
}
