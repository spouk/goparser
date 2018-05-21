package library

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	baseGoogle = "www.google.ru"
)

type (
	//base struct
	Google struct {
		BaseURL string
		Type    GoogleType
		Period  GoogleTime
		Params  GoogleParams
	}
	GoogleType struct {
		Image        string
		Video        string
		Applications string
		News         string
		Patents      string
		Books        string
	}
	GoogleTime struct {
		AnyTime   string
		Second    string
		Minute    string
		Minute10  string
		Hour      string
		Hour12    string
		Day       string
		Week      string
		Month     string
		Year      string
		TimeRange TimeRangeType
	}
	TimeRangeType struct {
		StartDate string
		EndDate   string
	}
	//extern structs
	GoogleParams struct {
		ImageParam  GoogleImageParam
		VideoParam  GoogleVideoParams
		RandomParam GoogleRandomParam
	}
	//RANDOM PARAMS
	GoogleRandomParam struct {
		CountResult CountResult
	}
	CountResult struct {
		Count50   string
		Count100  string
		Count500  string
		Count1000 string
	}
	//IMAGE TYPE
	GoogleImageParam struct {
		Prefix    string
		Size      ImageSize
		Color     ImageColor
		ImageType ImageType
	}
	ImageType struct {
		Face        string
		Photo       string
		Clipart     string
		LineDrawing string
		Animated    string
	}
	ImageSize struct {
		Large  string
		Medium string
		Icon   string
	}
	ImageColor struct {
		FullColor     string
		BlackAndWhite string
		Random        string
	}
	//VIDEO TYPE
	GoogleVideoParams struct {
		Prefix   string
		Duration VideoDuration
		Quality  VideoQuality
	}
	VideoDuration struct {
		Short  string
		Medium string
		Long   string
	}
	VideoQuality struct {
		High string
	}
)

func NewGoogleRandomParams() GoogleRandomParam {
	return GoogleRandomParam{
		CountResult: CountResult{
			Count50:   "50",
			Count100:  "100",
			Count500:  "500",
			Count1000: "1000",
		},
	}
}
func NewGoogleVideoParams() GoogleVideoParams {
	return GoogleVideoParams{
		Prefix: "tbs=",
		Duration: VideoDuration{
			Short:  "dur:s",
			Medium: "dur:m",
			Long:   "dur:l",
		},
		Quality: VideoQuality{
			High: "hq:h",
		},
	}
}
func NewImageType() ImageType {
	return ImageType{
		Face:        "itp:face",
		Photo:       "itp:photo",
		Clipart:     "itp:clipart",
		LineDrawing: "itp:lineart",
		Animated:    "itp:animated",
	}
}
func NewImageSize() ImageSize {
	return ImageSize{
		Large:  "isz:l",
		Medium: "isz:m",
		Icon:   "isz:i",
	}

}
func NewImageColor() ImageColor {
	return ImageColor{
		FullColor:     "ic:color",
		BlackAndWhite: "ic:gray",
		Random:        "ic:specific,isc:red,orange,yellow, green, teal, blue, purple, pink, white, gray, black, brown",
	}
}
func NewGoogleImageParam() GoogleImageParam {
	return GoogleImageParam{
		Prefix:    "tbs=",
		Size:      NewImageSize(),
		Color:     NewImageColor(),
		ImageType: NewImageType(),
	}
}
func NewGoogleParams() GoogleParams {
	return GoogleParams{
		ImageParam:  NewGoogleImageParam(),
		VideoParam:  NewGoogleVideoParams(),
		RandomParam: NewGoogleRandomParams(),
	}
}

//---------------------------------------------------------------------------
//  CONSTRUCTORS
//---------------------------------------------------------------------------
func NewGoogleType() GoogleType {
	return GoogleType{
		Image:        "tbm=isch",
		Video:        "tbm=vid",
		Applications: "tbm=app",
		News:         "tbm=nws",
		Patents:      "tbm=pts",
		Books:        "tbm=bks",
	}
}
func NewTimeRangeType(startDate, endDate string) TimeRangeType {
	return TimeRangeType{
		StartDate: startDate,
		EndDate:   endDate,
	}
}
func NewGoogleTime() GoogleTime {
	return GoogleTime{
		AnyTime:  "tbs=qdr:a",
		Second:   "tbs=qdr:s",
		Minute:   "tbs=qdr:n",
		Minute10: "tbs=qdr:n10",
		Hour:     "tbs=qdr:h",
		Hour12:   " tbs=qdr:h12",
		Day:      "tbs=qdr:d",
		Week:     "tbs=qdr:w",
		Month:    "tbs=qdr:m",
		Year:     "tbs=qdr:y",
	}
}
func NewGoogle(baseurl string) *Google {
	return &Google{
		BaseURL: baseurl,
		Type:    NewGoogleType(),
		Period:  NewGoogleTime(),
		Params:  NewGoogleParams(),
	}
}

//---------------------------------------------------------------------------
//  functions
//---------------------------------------------------------------------------
func (g *Google) NewSearch(request string, countResult string, typeReq string, period string, params ...string) (string) {
	//"http://www.google.com/search?q=%22michael+jackson%22&tbm=isch&tbs=ic:color,isz:lt,islt:4mp,itp:face,isg:to"
	r := fmt.Sprintf("http://%s/search?num=%s&q=%s&%s&%s%s",
		baseGoogle,
		countResult,
		url.PathEscape(request),
		typeReq,
		g.Params.ImageParam.Prefix,
		strings.Join(params, ","),
	)
	fmt.Printf("NewSearch: %v\n", r)
	return r
}
func (g *Google) makeRequest() {

}
