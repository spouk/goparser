package main

import (
	"time"
	"fmt"
	"errors"
	"log"
	"os"
	"net/http"
	"golang.org/x/net/html"
	"bytes"
	"yuristlic.parser/library"
	"net/url"

	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	//"crypto/tls"
	"io"

	"crypto/tls"
)

const (
	const_loggerName     = "[ yuristlic.parser ][ logger ] "
	const_version_parser = "0.0.1"
)

var (
	video  = "https://www.google.ru/search?q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%96%D0%9A%D0%A5+%D0%B0%D0%B2%D0%B0%D1%80%D0%B8%D0%B8+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&safe=off&biw=1096&bih=548&tbm=vid&source=lnms&sa=X&ved=0ahUKEwjQoNGx2ajWAhVCJpoKHch1CwwQ_AUICygC"
	images = "https://www.google.ru/search?safe=off&biw=1096&bih=548&tbm=isch&sa=1&q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%96%D0%9A%D0%A5+%D0%B0%D0%B2%D0%B0%D1%80%D0%B8%D0%B8+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&oq=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%96%D0%9A%D0%A5+%D0%B0%D0%B2%D0%B0%D1%80%D0%B8%D0%B8+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&gs_l=psy-ab.3...9594.45079.0.45343.41.38.2.0.0.0.114.2840.34j2.38.0....0...1.1.64.psy-ab..2.15.1117.0..0j0i67k1j0i30k1j0i8i30k1.66.a0f7Ae-MqgU"
	//"http://www.google.com/search?q=%22michael+jackson%22&tbm=isch&tbs=ic:color,isz:lt,islt:4mp,itp:face,isg:to"
	//http://www.google.ru/search?q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B%20%D0%B2%20%D0%96%D0%9A%D0%A5%20%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&tbm=isch&tbs=
	error_target = "ошибка: ошибка в названии ресурса или ресурс не доступен для обработки"
	b            = []byte(`<!DOCTYPE html>
<html>
    <head>
        <title>
            Title of the document
        </title>
    </head>
    <body>
    	<a href="http://spouk.ru" class="cls_1 cls_2 cls_3" id="idspoukru" title="simple title for href spoukru">LINK TO SPOUK.RU </a>
        body content
        <p>more content</p>
        <div class="panel-body">
        <form action="/user/login" method="post">

            <!--email-->
            <div class="panel-item">
                <label for="Robot">Почтовый адрес</label>
                <input type="email" id="Email" name="Email" placeholder="{{$stock.Desc.Email}}"
                       value="{{if $form.Email}} {{$form.Email}} {{end}}">
            </div>


            <div class="s-user-form-row">
                <label for="Robot">Почтовый адрес</label>
                <input type="email" id="Email" name="Email" placeholder="{{$stock.Desc.Email}}"
                       value="{{if $form.Email}} {{$form.Email}} {{end}}">
                <br>
            </div>
            <div class="s-user-form-row">
                <label for="ErrorEmail"></label>
                <i class="s-user-error " id="ErrorEmail">
                    {{if $stock.Error.Email.Error}}
                    <span>{{$stock.Error.Email.Error}}</span>
                    {{end}}
                </i>
            </div>


            <div class="s-user-form-row">
                <label for="Robot">Пароль</label>
                <input type="text" id="Password" name="Password" placeholder="{{$stock.Desc.Password}}"
                       value="{{if $form.Password}} {{$form.Password}} {{end}}"> <br>
            </div>
            <div class="s-user-form-row">
                <label for="ErrorPassword"></label>
                <i class="s-user-error " id="ErrorPassword">
                    {{if $stock.Error.Password.Error}}
                    <span>{{$stock.Error.Password.Error}}</span>
                    {{end}}
                </i>
            </div>

            <div class="s-user-form-row">
                <label for="Robot">Я не робот</label>
                <input type="checkbox" name="Robot" id="Robot">
            </div>
            <div class="s-user-form-row">
                <label for="ErrorRobot"></label>
                <i class="s-user-error " id="ErrorRobot">
                    {{if $stock.Error.Robot.Error}}
                    <span>{{$stock.Error.Robot.Error}}</span>
                    {{end}}
                </i>
            </div>
            <div class="s-user-button">
                <button type="submit">Пройти авторизацию</button>
            </div>
        </form>
    </div>
    </body>
</html>`)
)

type (
	ParserExample struct {
		//public methods&data
		Version string
		//private methods&data
		stockResult chan resultWorker
		stockErrors chan errorWorker
		//logging
		logger *log.Logger
		//headers
		headers library.Headers
	}
	resultWorker struct {
		NameWorker   string
		TargerWorker string
		ObjectResult interface{}
	}
	errorWorker struct {
		NameWorker   string
		TargerWorker string
		ErrorWorker  error
	}
	worker struct {
		Name       string
		TimeStart  time.Time
		TimeEnd    time.Time
		resultChan *chan resultWorker
	}
)

func main() {
	ns := library.NewStack(
		"ЖКХ прикольные картинки",
		library.TypeRequest{Image: true},
		library.TypeSearchEngine{Yandex: true},
		10,
	)
	ns.RandomGenerate(100, library.TypeRequest{Image: true}, "ЖКХ Прикольные картинки", library.TypeSearchEngine{Yandex: true})
	ns.RandomGenerate(200, library.TypeRequest{Image: true}, "ЖКХ Прикольные картинки2", library.TypeSearchEngine{Yandex: true})
	ns.RandomGenerate(300, library.TypeRequest{Image: true}, "ЖКХ Прикольные картинки3", library.TypeSearchEngine{Yandex: true})
	ns.RandomGenerate(400, library.TypeRequest{Image: true}, "ЖКХ Прикольные картинки4", library.TypeSearchEngine{Yandex: true})
	ns.RandomGenerate(500, library.TypeRequest{Image: true}, "ЖКХ Прикольные картинки5", library.TypeSearchEngine{Yandex: true})
	ns.ShowStock()
	ns.Run(10)
	fmt.Printf("Count stock: %v\n", len(ns.Stock))
	for k, v := range ns.Stock {
		fmt.Printf("%v : %v\n", k, len(v))
	}

	os.Exit(1)

	//u := "https://gobyexample.com/worker-pools"
	//u2 := "https://www.google.ru/search?q=golang+html+parse+example&oq=golang+html+parse+example&aqs=chrome..69i57j0j69i64.3758j0j7&sourceid=chrome&ie=UTF-8"
	p := NewParser()
	fmt.Printf("Hello %v\n", p.NewError(error_target))
	p.logger.Printf("Hello world")
	//p.HttpRequest(u)
	//p.HttpRequest(u2)
	tt := html.NewTokenizer(bytes.NewBuffer(b))
	fmt.Printf("Tokeninize: %v\n", tt)
	for {
		t := tt.Next()
		if t == html.ErrorToken {
			fmt.Printf("Found end buffer/page")
			break
		}
		if t == html.StartTagToken {
			r := tt.Token()
			fmt.Printf("\n-------------------------------------------\nNewToken\nData: %v\nAttr: %v\nType: %v\n", r.Data, r.Attr, r.Type)

			//type attributes
			for _, x := range r.Attr {
				fmt.Printf("Key: `%s` %20s Value: `%20s`\n", x.Key, " ", x.Val)
			}
		}
	}
	//---------------------------------------------------------------------------
	//  googlerequest tsting
	//---------------------------------------------------------------------------
	base := "www.google.ru"
	r := library.NewGoogle(base)
	fmt.Printf("\nGoogleRequest: %v\b", r)
	fmt.Printf("\nQuery escape: %v\n", url.PathEscape("пример запроса"))
	//example IMAGE
	resultSearchImages := r.NewSearch("приколы в ЖКХ картинки",
		r.Params.RandomParam.CountResult.Count500,
		r.Type.Image,
		r.Period.AnyTime,
		r.Params.ImageParam.ImageType.Animated,
		r.Params.ImageParam.ImageType.Photo,
		r.Params.ImageParam.Size.Large,
		r.Params.ImageParam.Size.Medium,
		r.Params.ImageParam.Color.FullColor,
	)
	//example VIDEO
	resultSearchVideo := r.NewSearch("Приколы в ЖКХ видео",
		r.Params.RandomParam.CountResult.Count500,
		r.Type.Video,
		r.Period.AnyTime,
		r.Params.VideoParam.Duration.Short,
		r.Params.VideoParam.Quality.High,
	)
	fmt.Printf("ResultVideo: %v\nResultImages: %v\n", resultSearchVideo, resultSearchImages)
	//d1 := "http://www.google.ru/search?q=%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%B2+%D0%96%D0%9A%D0%A5+%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&newwindow=1&ie=UTF-8&tbs=itp:photo,ic:color,isz:i&tbm=isch&source=lnt&sa=X&ved=0ahUKEwihkfzx66vWAhVHQZoKHbT0CewQpwUIDg"
	//d2 := "http://spouk.ru/searhc?q=simple+tester"
	//neweq := "https://www.google.ru/search?num=100&safe=off&hl=ru&q=%D0%9F%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%B2+%D0%96%D0%9A%D0%A5+%D0%B2%D0%B8%D0%B4%D0%B5%D0%BE&oq=%D0%9F%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B+%D0%B2+%D0%96%D0%9A%D0%A5+%D0%B2%D0%B8%D0%B4%D0%B5%D0%BE"
	//neweq := "https://yandex.ru/search/?text=%D0%96%D0%9A%D0%A5%20%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B%20%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&lr=10834"
	neweq := "https://yandex.ru/images/search?text=%D0%96%D0%9A%D0%A5%20%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B%20%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8"
	//deep = https://yandex.ru/images/search?p=4&text=%D0%B6%D0%BA%D1%85%20%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B&rpt=image

	//p.HttpRequest(neweq)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://ya.ru/?text=sex+images", nil)
	req.Header.Add("UserAgent", "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; .NET CLR 2.0.50727; .NET CLR 1.0.3705; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729; .NET CLR 1.1.4322; .NET4.0C; .NET4.0E; IPH 1.1.21.4019; chromeframe/24.0.1312.52)")
	fmt.Printf("REQ: %v\nClient:%v\n", req, client)
	resp, _ := client.Do(req)
	rrrr, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("RESULT BODY: %s\v", string(rrrr), neweq)

	//u, _ := url.Parse(d2)
	//v, _ := url.ParseQuery(d2)

	//fmt.Printf("U: %v\n", u.Path, u.Host, u.Opaque, u.RawPath)
	//fmt.Printf("Values: %v\n", u.Query().Get("q"))
	//for key, value := range v {
	//	fmt.Printf("Key: %v    Value: %v\n", key, value)
	//}
	//rand.Seed(time.Now().UTC().UnixNano())
	//fmt.Printf("Random: %v\n", rand.Intn(4))

	//---------------------------------------------------------------------------
	//  EXMAPLE GO QUERY
	//---------------------------------------------------------------------------

	//ExampleScrape("https://yandex.ru/search/?text=ЖКХ")
	//req1 := "https://yandex.ru/search/?text=porno"
	//req2 := "https://google.ru/search?q=porno"
	//req3 := "https://www.google.ru/search?q=porno&num=100&safe=off&source=lnms&tbm=isch"

	req4 := "https://yandex.ru/images/search?text=%D0%96%D0%9A%D0%A5%20%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B%20%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8"

	//links

	//more
	//https://yandex.ru/images/search?p=5&nomisspell=1&text=%D0%B6%D0%BA%D1%85%20%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&source=related-8&rpt=image&uinfo=sw-1920-sh-1080-ww-1744-wh-872-pd-1.100000023841858-wp-16x9_2560x1440

	//more
	//https://yandex.ru/images/search?p=10&nomisspell=1&text=%D0%B6%D0%BA%D1%85%20%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&source=related-8&rpt=image&uinfo=sw-1920-sh-1080-ww-1744-wh-872-pd-1.100000023841858-wp-16x9_2560x1440

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
	resp, err := client.Get(req4)
	if err != nil {
		fmt.Println(err)
	}
	//f, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("RESULT YANEX SEARCH: %v\n", string(f))
	ShowAllTags(resp.Body)

	//---------------------------------------------------------------------------
	//  PARSE ANSWER FROM YANDEX
	//---------------------------------------------------------------------------
	href := "/images/search?text=%D0%96%D0%9A%D0%A5%20%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B%20%D0%BA%D0%B0%D1%80%D1%82%D0%B8%D0%BD%D0%BA%D0%B8&amp;img_url=https%3A%2F%2Fcs9.pikabu.ru%2Fpost_img%2F2017%2F04%2F15%2F11%2F14922805491814924.jpg&amp;pos=2&amp;rpt=simage"
	uu, _ := url.Parse(href)
	fmt.Printf("Text :%v\n", uu.Query().Get("text"))
	fmt.Printf("Text :%v\n", uu.Query().Get("img_url"))

}
func ShowAllTags(resp io.Reader) bool {
	t := html.NewTokenizer(resp)
	for {
		tt := t.Next()
		switch {
		case tt == html.ErrorToken:
			return false
		case tt == html.StartTagToken:
			t1 := t.Token()
			linkA := t1.Data == "a"
			if linkA {
				//fmt.Printf("Found link: %v\nAttrs: %v\n", t1.Data, t1.Attr)

				for _, a := range t1.Attr {
					if a.Key == "href" {
						fmt.Printf("Value: %v\n", a.Val)
						u, err := url.Parse(a.Val)
						if err != nil {
							fmt.Printf("ERROR: %v\n", err.Error())
						} else {
							fmt.Printf("===> %s\n", u.Query().Get("text"))
							fmt.Printf("===> %s\n", u.Query().Get("img_url"))
						}
						//fmt.Printf("URL UMAGES: %v\n", u.Query().Get("img_url"))
						//fmt.Printf("[%v] Href: %v\n", t1.Data, a.Val)

						break
					}

				}
			}
		}
	}
}
func ExampleScrape(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Printf("LOGFATAL: %v\n", err)
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		band := s.Find("a").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}

func NewParser() *ParserExample {
	p := &ParserExample{
		Version:     const_version_parser,
		stockErrors: make(chan errorWorker),
		stockResult: make(chan resultWorker),
		logger:      log.New(os.Stdout, const_loggerName, log.Llongfile|log.Ldate|log.Ltime),
		headers:     library.Headers{},
	}
	return p
}
func (p *ParserExample) NewError(msg string) error {
	return errors.New(fmt.Sprintf("[%s][ошибка] `%s`", const_loggerName, msg))
}

//---------------------------------------------------------------------------
//  part http request
//---------------------------------------------------------------------------
func (p *ParserExample) HttpRequest(url string) (bool, error) {
	//make client + add headers for validate search requests
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		p.logger.Printf(err.Error())
		return false, err
	}
	p.headers.AssignHeader(req)
	p.logger.Printf("Client: %v\n", client)

	resp, err := client.Do(req)
	if err != nil {
		p.logger.Printf(err.Error())
		return false, err
	}
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("RESP BODY: %v\n", string(b))
	return false, nil

	//resp, err := http.Get(url)
	//if err != nil {
	//	p.logger.Printf(err.Error())
	//	return false, err
	//}
	//bytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	p.logger.Printf(err.Error())
	//	return false, err
	//}
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
				p.logger.Printf("Found link: %v\nAttrs: %v\n", t1.Data, t1.Attr)

				for _, a := range t1.Attr {
					if a.Key == "href" {
						p.logger.Printf("[%v] Href: %v\n", t1.Data, a.Val)
						break
					}

				}
			}
		}
	}
}
