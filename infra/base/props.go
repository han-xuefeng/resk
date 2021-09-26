package base

import (
	"github.com/tietang/props/kvs"
	"study-gin/resk/infra"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter)Init(content infra.StarterContent) {
	// 初始化配置  在bootStarter实例化的时候就已经完成了，这里是验证有没有加载配置文件
	props = content.Props() //加到全局变量里面
	//fmt.Println("配置文件加载")
}