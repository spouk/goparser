package library


const (
	YANDEXBASE  = "https://yandex.ru"
	YANDEXBASEIMAGE = "https://yandex.ru/images/search?p=%d&text=%s"
	YANDEXBASEVIDEO = "https://yandex.ru/video/search?p=%d&text=%s"


	//YANDEXBASEIMAGE = "https://yandex.ru/images/search?p=1&text=&rpt=image"
	//"https://yandex.ru/images/search?p=1&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B"
	//https://yandex.ru/video/search?p=2&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
	//https://yandex.ru/video/search?p=3&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
	//https://yandex.ru/video/search?p=4&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
	//https://yandex.ru/video/search?p=4&text=%D0%B6%D0%BA%D1%85+%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B
//<a class="serp-item__link" href="/images/search?p=4&amp;text=%D0%B6%D0%BA%D1%85%20%D0%BF%D1%80%D0%B8%D0%BA%D0%BE%D0%BB%D1%8B&amp;img_url=http%3A%2F%2Fs02.radikal.ru%2Fi175%2F1707%2F24%2F97a453cd5936.jpg&amp;pos=161&amp;rpt=simage"><img class="serp-item__thumb" src="//im0-tub-ru.yandex.net/i?id=94b4bca56887f696bb73f0561471e60e&amp;n=13" onerror="this.onerror = &quot;&quot;;var item = this.parentNode.parentNode.parentNode;item.parentNode.removeChild(item);window.Ya.Images.errcnt(this, 71729);return true;" alt="Приколы омского ЖКХ." style="height: 181px;"><div class="serp-item__plates"><div class="serp-item__meta">960×721</div></div></a>
)
type (
	YandexParser struct {
		CountPages int
	}
)

func NewYandexParser(countPages int) *YandexParser{
	return &YandexParser{
		CountPages:countPages,
	}
}
func(y *YandexParser) SearchImages() {

}
func(y *YandexParser) SearchVideo() {

}