package main

import (
	"fmt"
	"strings"

	//"siteGhost3/cmd/metaBuilder/conf"
	"siteGhost3/cmd/metaBuilder/spider"
	"siteGhost3/conf"
	"siteGhost3/model"
	"siteGhost3/utils"

	"github.com/BurntSushi/toml"
	log "github.com/aiyaya188/go-libs/logger"
)

func testTrim() {
	str := `<div class="answer-content js-resolved-answer">双子 　　女：88% 　　此类女生天生丽质，通常是炯炯有神的双眼、嫩白的<a data-id="link-to-so" href="http://www.so.com/s?q=%E7%9A%AE%E8%82%A4&amp;ie=utf-8&amp;src=internal_wenda_recommend_textn" style="color:#0063c8;cursor:pointer;" text="详情页文字内链_点击实体词" target="_blank">皮肤</a>，小巧的个子。 　　 巨蟹 　　女：75% 　　此类女生温柔娴淑、属于大家闺秀型，当然也很漂亮了哦！ 　　 狮子 　　女：70% 　　此类女生和双子一样，属于天生丽质型，再加上她贵族的<a data-id="link-to-so" href="http://www.so.com/s?q=%E6%B0%94%E8%B4%A8&amp;ie=utf-8&amp;src=internal_wenda_recommend_textn" style="color:#0063c8;cursor:pointer;" text="详情页文字内链_点击实体词" target="_blank">气质</a>~让人连连叫好处女 　　女：50% 　　此类女生有的漂亮，有的普通，但他们笑起来都是很迷人的天枰座 　　女：89% 　　此类女生在<a data-id="link-to-so" href="http://www.so.com/s?q=%E4%BB%BB%E4%BD%95%E5%9C%B0%E6%96%B9&amp;ie=utf-8&amp;src=internal_wenda_recommend_textn" style="color:#0063c8;cursor:pointer;" text="详情页文字内链_点击实体词" target="_blank">任何地方</a>都是大家闺秀，<a data-id="link-to-so" href="http://www.so.com/s?q=%E9%95%BF%E7%9B%B8&amp;ie=utf-8&amp;src=internal_wenda_recommend_textn" style="color:#0063c8;cursor:pointer;" text="详情页文字内链_点击实体词" target="_blank">长相</a>得体，非常漂亮！ 　　 天蝎座 　　女：75% 　　此类女生多数<a data-id="link-to-so" href="http://www.so.com/s?q=%E6%80%A7%E6%84%9F&amp;ie=utf-8&amp;src=internal_wenda_recommend_textn" style="color:#0063c8;cursor:pointer;" text="详情页文字内链_点击实体词" target="_blank">性感</a>，让人觉得很有吸引力 　　 射手座 　　女：80% 　　此类女生活泼开朗、可爱温柔。自然让人觉得很漂亮 　　 摩羯座 　　女：60% 　　此类女生大多数较普通，但很有气质 　　 水瓶座 　　女：70% 　　此类女生很可爱，笑起来更是这样。 　　 双鱼座 　　女：85% 　　此类女生是典型的小家碧玉，通常都有着楚楚可怜的大眼睛。 　　 白羊座 　　女：漂亮指数：40% 　　此类女生可能胖了点，又可能瘦了点，但本质是漂亮的^_^ 　　 金牛 　　女：55% 　　此类女生可能长得知算普通，但有特殊的吸引力。只是参考</div>`
	str = utils.TrimHtml(str)
	fmt.Println("str:", strings.TrimSpace(str))
}
func main() {
	var err error
	_, err = toml.DecodeFile("config.toml", &conf.Config)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.InitLogger(conf.Config.Logging)
	var wd spider.Wenda360
	wd.Host = "https://wenda.so.com/"
	data := wd.SpiderContent("https://wenda.so.com/q/1378840745065566")
	fmt.Println("data:", data)
	return

	model.InitDB(conf.Config, 1)
	model.Migration()
	utils.CreateHttpClient()
	model.Migration()
	server := Server{}
	server.run()

}
