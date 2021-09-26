package infra

func init() {
	Register(&ConfStarter{})
}

type ConfStarter struct {
	BaseStarter
}

func (c *ConfStarter)Init(StarterContent)  {
	println("配置初始化")
}

func (c *ConfStarter)Setup(StarterContent)  {
	println("配置安装")
}

func (c *ConfStarter)Start(StarterContent)  {
	println("配置启动")
}

