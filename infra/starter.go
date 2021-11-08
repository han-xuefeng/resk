package infra

import (
	"github.com/tietang/props/kvs"
	"sort"
)

const (
	KeyProps = "_conf"
)

// 基础资源上下文结构体
type StarterContext map[string]interface{}

func (s StarterContext)Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置还没有初始化")
	}
	return p.(kvs.ConfigSource)
}

func (s StarterContext) SetProps(conf kvs.ConfigSource) {
	s[KeyProps] = conf
}

//基础资源启动器接口
type Starter interface {
	//1. 系统启动的时候，初始化一些基础资源
	Init(StarterContext)
	//2. 系统基础资源的安装
	Setup(StarterContext)
	//3. 启动系统基础资源
	Start(StarterContext)

	//启动器是否阻塞
	StartBlocking() bool

	//4. 关闭资源
	Stop(StarterContext)
	PriorityGroup() PriorityGroup
	Priority() int
}

// 验证基础类是否实现了接口的所有方法
var _ Starter = new(BaseStarter)

type PriorityGroup int

const (
	SystemGroup         PriorityGroup = 30
	BasicResourcesGroup PriorityGroup = 20
	AppGroup            PriorityGroup = 10

	INT_MAX          = int(^uint(0) >> 1)
	DEFAULT_PRIORITY = 10000
)

//基础类   空实现
type BaseStarter struct {

}

func (b *BaseStarter) Init(StarterContext) {

}
func (b *BaseStarter) Setup(StarterContext) {

}
func (b *BaseStarter) Start(StarterContext) {

}
func (b *BaseStarter)StartBlocking() bool {
	return false
}
func (b *BaseStarter) Stop(StarterContext) {

}
func (b *BaseStarter) PriorityGroup() PriorityGroup { return BasicResourcesGroup }
func (b *BaseStarter) Priority() int                { return DEFAULT_PRIORITY }

type Starters []Starter

func (s Starters) Len() int      { return len(s) }
func (s Starters) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Starters) Less(i, j int) bool {
	return s[i].PriorityGroup() > s[j].PriorityGroup() && s[i].Priority() > s[j].Priority()
}

// 启动器注册器
type starterRegister struct {
	Starters []Starter
}

func SortStarters() {
	sort.Sort(Starters(StarterRegister.AllStarters()))
}

//获取所有注册的starter
func GetStarters() []Starter {
	return StarterRegister.AllStarters()
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

