package model

import (
	ds "github.com/aiyaya188/go-libs/ds"
	"github.com/jinzhu/gorm"
)

const (
	URLWaitSpider = 0 //地址等待抓取
	URLFinish     = 1 //地址抓取完毕
)

type UrlData struct {
	gorm.Model
	Title    string
	Url      string `gorm:"unique_index:unique_index_url"`
	Status   int
	KeyWord  string
	Res      string
	Resource string //来源
}

func CreateUrl(keyword, title, url, resource string) error {
	var ud UrlData
	ud.Title = title
	ud.Url = url
	ud.KeyWord = keyword
	ud.Resource = resource
	if err := ds.DB.Table("url_data").Create(&ud).Error; err != nil {
		return err
	}
	return nil
}
