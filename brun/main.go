package main

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	_ "study-gin/resk"
	"study-gin/resk/infra"
	"study-gin/resk/infra/base"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
}
