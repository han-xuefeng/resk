package infra

import (
	"github.com/tietang/props/kvs"
)

type BootApplication struct {
	conf           kvs.ConfigSource
	starterContent StarterContext
}


func New (conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		conf:           conf,
		starterContent: StarterContext{},
	}
	b.starterContent[KeyProps] = conf
	return b
}

func (b *BootApplication) Start() {
	//1. 初始化所有的Start
	b.init()
	//2. 安装
	b.setup()
	//3. 启动
	b.start()
}


func (b *BootApplication) init() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(b.starterContent)
	}
}

func (b *BootApplication) setup() {
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(b.starterContent)
	}
}

func (b *BootApplication) start() {
	for i, starter := range StarterRegister.AllStarters() {
		if starter.StartBlocking() && i != len(StarterRegister.AllStarters()) -1 {
			go starter.Start(b.starterContent)
		} else {
			starter.Start(b.starterContent)
		}
	}
}
