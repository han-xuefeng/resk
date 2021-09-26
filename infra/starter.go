package infra

import "github.com/tietang/props/kvs"

const (
	KeyProps = "_conf"
)

// 基础资源上下文结构体
type StarterContent map[string]interface{}

func (s StarterContent)Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没有初始化")
	}
	return p.(kvs.ConfigSource)
}

//基础资源启动器接口
type Starter interface {
	//1. 系统启动的时候，初始化一些基础资源
	Init(StarterContent)
	//2. 系统基础资源的安装
	Setup(StarterContent)
	//3. 启动系统基础资源
	Start(StarterContent)

	//启动器是否阻塞
	StartBlocking() bool

	//4. 关闭资源
	Stop(StarterContent)
}

// 验证基础类是否实现了接口的所有方法
var _ Starter = new(BaseStarter)

//基础类   空实现
type BaseStarter struct {

}

func (b *BaseStarter) Init(StarterContent) {

}
func (b *BaseStarter) Setup(StarterContent) {

}
func (b *BaseStarter) Start(StarterContent) {

}
func (b *BaseStarter)StartBlocking() bool {
	return false
}
func (b *BaseStarter) Stop(StarterContent) {

}

// 启动器注册器
type starterRegister struct {
	Starters []Starter
}

func (r *starterRegister) Register(s Starter) {
	r.Starters = append(r.Starters, s)
}

func (r *starterRegister) AllStarters() []Starter {
	return r.Starters
}

var StarterRegister *starterRegister = new(starterRegister)

func Register(s Starter) {
	StarterRegister.Register(s)
}

// 系统资源的启动
func SystemRun() {
	ctx := StarterContent{}
	for _, starter := range StarterRegister.AllStarters() {
		starter.Init(ctx)
	}
	for _, starter := range StarterRegister.AllStarters() {
		starter.Setup(ctx)
	}
	for _, starter := range StarterRegister.AllStarters() {
		starter.Start(ctx)
	}
}

