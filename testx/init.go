package testx

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"study-gin/resk/infra"
	"study-gin/resk/infra/base"
)

func init(){
	//获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("../brun/config.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)

	//

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()
}