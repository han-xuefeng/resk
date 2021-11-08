package infra

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"reflect"
)

type BootApplication struct {
	IsTest     bool
	conf           kvs.ConfigSource
	starterCtx StarterContext
}


func New (conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		conf:           conf,
		starterCtx: StarterContext{},
	}
	b.starterCtx.SetProps(conf)
	return b
}

func (e *BootApplication) Start() {
	//1. 初始化所有的Start
	e.init()
	//2. 安装
	e.setup()
	//3. 启动
	e.start()
}


//程序初始化
func (e *BootApplication) init() {
	log.Info("Initializing starters...")
	for _, v := range GetStarters() {
		typ := reflect.TypeOf(v)
		log.Debugf("Initializing: PriorityGroup=%d,Priority=%d,type=%s", v.PriorityGroup(), v.Priority(), typ.String())
		v.Init(e.starterCtx)
	}
}

//程序安装
func (e *BootApplication) setup() {

	log.Info("Setup starters...")
	for _, v := range GetStarters() {
		typ := reflect.TypeOf(v)
		log.Debug("Setup: ", typ.String())
		v.Setup(e.starterCtx)
	}

}

//程序开始运行，开始接受调用
func (e *BootApplication) start() {

	log.Info("Starting starters...")
	for i, v := range GetStarters() {
		typ := reflect.TypeOf(v)
		log.Debug("Starting: ", typ.String())
		if v.StartBlocking() {
			if i+1 == len(GetStarters()) {
				v.Start(e.starterCtx)
			} else {
				go v.Start(e.starterCtx)
			}
		} else {
			v.Start(e.starterCtx)
		}

	}
}

//程序开始运行，开始接受调用
func (e *BootApplication) Stop() {

	log.Info("Stoping starters...")
	for _, v := range GetStarters() {
		typ := reflect.TypeOf(v)
		log.Debug("Stoping: ", typ.String())
		v.Stop(e.starterCtx)
	}
}
