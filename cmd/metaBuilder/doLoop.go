package main

import (
	"context"
	"siteGhost3/cmd/metaBuilder/spider"
	"siteGhost3/model"
	"time"

	"golang.org/x/sync/errgroup"
)

func GetKeyword() string {
	keyword := "星座"
	return keyword
}

func CreateURL(data []spider.URLData, key string, res string) {
	for _, v := range data {
		model.CreateUrl(key, v.Title, v.URL, res)
	}

}
func DoLoopJob() {
	wd360Host := "https://wenda.so.com"
	wd360 := spider.Wenda360{Host: wd360Host}

	wg, _ := errgroup.WithContext(context.Background())
	for {
		keyword := GetKeyword()
		wg.Go(func() error {
			data := wd360.SpiderURL(keyword)
			CreateURL(data, keyword, "360")
			return nil
		})
		wg.Wait()
		time.Sleep(3 * time.Second)
	}
}
