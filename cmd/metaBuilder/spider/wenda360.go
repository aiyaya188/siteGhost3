package spider

import (
	"context"
	"log"
	"strings"

	"siteGhost3/utils"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Wenda360 struct {
	Host string
}

func (w *Wenda360) SpiderURL(keyword string) []URLData {
	var content string
	var urlData []URLData
	ua := browser.Computer()
	//fmt.Println("ua:", ua)
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(ua),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	//创建chrome窗口
	err := chromedp.Run(ctx,
		//chromedp.Navigate(`https://wenda.so.com/`),
		chromedp.Navigate(w.Host),
		chromedp.WaitVisible(`div#ft`),
		chromedp.SendKeys(`input#js-sh-ipt`, keyword, chromedp.NodeVisible),
		chromedp.Click(`input.js-suggest-search-btn`, chromedp.NodeVisible),
		chromedp.WaitVisible(`div#footer`),
		chromedp.OuterHTML(`ul#js-qa-list`, &content, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("content:%s", content)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("div>h3>a").Text()
		addr, _ := selection.Find("div>h3>a").Attr("href")
		if addr != "" && s != "" {
			urlData = append(urlData, URLData{URL: w.Host + addr, Title: s})
		}
	})
	//fmt.Println("点击")
	err = chromedp.Run(ctx,
		chromedp.Click(`a.next`, chromedp.NodeVisible),
		chromedp.OuterHTML(`ul#js-qa-list`, &content, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}
	dom, err = goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatalln(err)
	}
	dom.Find("li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find("div>h3>a").Text()
		addr, _ := selection.Find("div>h3>a").Attr("href")
		if addr != "" && s != "" {
			urlData = append(urlData, URLData{URL: w.Host + addr, Title: s})
		}
	})
	return urlData
}

func (w *Wenda360) SpiderContent(url string) []string {
	var content string
	ua := browser.Computer()
	//fmt.Println("ua:", ua)
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("headless", false),
		chromedp.UserAgent(ua),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	//创建chrome窗口
	err := chromedp.Run(ctx,
		//chromedp.Navigate(`https://wenda.so.com/`),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div.main-nav`),
		chromedp.OuterHTML(`div#answer`, &content, chromedp.NodeVisible),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("content:%s", content)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatalln(err)
	}
	var answer []string
	dom.Find("div.answer-part>div.answer-content").Each(func(i int, selection *goquery.Selection) {
		s := selection.Text()
		log.Printf("s:%s", s)
		answer = append(answer, utils.TrimHtml(s))

	})
	return answer

}
