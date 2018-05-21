package lib

const (
	SEARCH_TYPE_YANDEX = "yandex"
	SEARCH_TYPE_GOOGLE = "google"

	SEARCH_TYPE_CONTEXT_VIDEO  = "video"
	SEARCH_TYPE_CONTEXT_IMAGES = "images"
)

type SearchRequest struct {
	RequestDeep   int
	RequestSearch string
	RequestType   string
	RequestText   string
	RequestURL    string
	RequestDesc   string
}
type SearchResult struct {
	SReq     *SearchRequest
	Getting  bool
	LinkFile string
	Name     string
	Ext      string
	Size     int64
	Type     string
	Time     int64
	Desc     string
}
