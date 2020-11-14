package model

import (
	"github.com/aiyaya188/go-libs/conf"
	ds "github.com/aiyaya188/go-libs/ds"
	log "github.com/aiyaya188/go-libs/logger"
)

//Migration 数据迁徙
func InitDB(cfg conf.Config, single int) {
	log.Infof("database:%s,type:%s", cfg.DataSource.URL, cfg.DataSource.Dialect)
	ds.InitDS(cfg)
	if single == 1 {
		ds.DB.SingularTable(true)
	}
	//migration()
	//	DataSeed()
}

func Migration() {
	var tables []interface{}
	tables = append(tables, UrlData{})
	log.Info("migration len:", len(tables))
	for _, v := range tables {
		ds.DB.AutoMigrate(v)
	}
}
