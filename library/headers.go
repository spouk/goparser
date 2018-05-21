package library

import (
	"time"
	"net/http"
	"math/rand"
)

type Headers struct {
	HeaderStock []RandomHeader
}
type RandomHeader map[string]string

func (h *Headers) AssignHeader(req *http.Request) error {
	for key, value := range  h.GetRandomHeader() {
		req.Header.Add(key, value)
	}
	return nil
}
func (h *Headers) AddHeader(n RandomHeader) error {
	h.HeaderStock = append(h.HeaderStock, n)
	return nil
}
func (h *Headers) GetRandomHeader() RandomHeader {
	rand.Seed(time.Now().UTC().UnixNano())
	randomIndex := rand.Intn(len(h.HeaderStock))
	return h.HeaderStock[randomIndex]
}
func NewHeaders() Headers {
	//mozzila Linux
	r1 := RandomHeader{}
	r1["UserAgent"] = "Mozilla/5.0 (X11; Linux x86_64; rv:10.0.12) Gecko/20100101 Firefox/10.0.12 Iceweasel/10.0.12"
	r1["Connection"] = "keep-alive"
	r1["AcceptEncooding"] = "deflate"
	r1["AcceptLanguage"] = "ru-ru,ru;q=0.8,en-us;q=0.5,en;q=0.3"
	r1["CasheControl"] = "max-age=0"
	r1["AcceptCharset"] = "utf-8;q=0.7,*;q=0.3"
	r1["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"

	//mozzila windows
	r2 := RandomHeader{}
	r2["UserAgent"] = "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:6.0a2) Gecko/20110613 Firefox/6.0a2"
	r2["Connection"] = "keep-alive"
	r2["AcceptEncooding"] = "deflate"
	r2["AcceptLanguage"] = "ru-ru,ru;q=0.8,en-us;q=0.5,en;q=0.3"
	r2["CasheControl"] = "max-age=0"
	r2["AcceptCharset"] = "utf-8;q=0.7,*;q=0.3"
	r2["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"

	//mozzila Iron
	r3 := RandomHeader{}
	r3["UserAgent"] = "Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2050.0 Iron/38.0.2150.0 Safari/537.36"
	r3["Connection"] = "keep-alive"
	r3["AcceptEncooding"] = "deflate"
	r3["AcceptLanguage"] = "ru-ru,ru;q=0.8,en-us;q=0.5,en;q=0.3"
	r3["CasheControl"] = "max-age=0"
	r3["AcceptCharset"] = "utf-8;q=0.7,*;q=0.3"
	r3["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"

	//Macintosh
	r4 := RandomHeader{}
	r4["UserAgent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_2) AppleWebKit/536.26.17 (KHTML, like Gecko) Version/6.0.2 Safari/536.26.17"
	r4["Connection"] = "keep-alive"
	r4["AcceptEncooding"] = "deflate"
	r4["AcceptLanguage"] = "ru-ru,ru;q=0.8,en-us;q=0.5,en;q=0.3"
	r4["CasheControl"] = "max-age=0"
	r4["AcceptCharset"] = "utf-8;q=0.7,*;q=0.3"
	r4["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"

	//Macintosh
	r5 := RandomHeader{}
	r5["UserAgent"] = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.77 Safari/537.36 Vivaldi/1.7.735.27"
	r5["Connection"] = "keep-alive"
	r5["AcceptEncooding"] = "deflate"
	r5["AcceptLanguage"] = "ru-ru,ru;q=0.8,en-us;q=0.5,en;q=0.3"
	r5["CasheControl"] = "max-age=0"
	r5["AcceptCharset"] = "utf-8;q=0.7,*;q=0.3"
	r5["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"

	return Headers{
		HeaderStock: []RandomHeader{
			r1, r2, r3, r4, r5,
		},
	}
}

