package main

import (
	"fmt"

	"github.com/aiyaya188/go-libs/surf"
)

/*
func main() {
	addr := "https://ip.ihuan.me/address/5Lit5Zu9.html" //国内
	//addr := "https://ip.ihuan.me/address/5Lit5Zu9.html"//国外
	browser := surf.NewBrowser()
	err := browser.Open(addr)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	ipCount := 0
	flag := true
	for flag {
		browser.Dom().Find("tbody").Each(func(_ int, s *goquery.Selection) {
			//fmt.Println(s.Text())
			fmt.Println("tr 长度:", s.Find("tr").Length())
			ipCount = s.Find("tr").Length()
			if ipCount < 3 {
				flag = false
				return
			}
			s.Find("tr").Each(func(_ int, tr *goquery.Selection) {
				fmt.Println("ip:", tr.Find("td:nth-child(1)").Text())
				fmt.Println("port:", tr.Find("td:nth-child(2)").Text())
				fmt.Println("location:", tr.Find("td:nth-child(3)").Text())
				fmt.Println("https:", tr.Find("td:nth-child(5)").Text())
				fmt.Println("post:", tr.Find("td:nth-child(6)").Text())
			})
		})
		time.Sleep(2 * time.Second)
		if err := browser.Click("a[aria-label=Next]"); err != nil {
			flag = false
		}
	}
}
*/
func main() {
	addr := "https://wenda.so.com" //国内
	browser := surf.NewBrowser()
	err := browser.Open(addr)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	input := browser.Dom().Find("#js-sh-ipt")

	fmt.Println("title:", browser.Title())
}
