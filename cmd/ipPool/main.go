package main

import (
	"fmt"
	"siteGhost3/cmd/ipPool/spider"
)

func main() {
	addr := "https://ip.ihuan.me/address/5Lit5Zu9.html"
	err := spider.GetTargetContent(addr)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

}
