package base

import (
	"github.com/tietang/props/kvs"
	"study-gin/resk/infra"
	"sync"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	return props
}

type PropsStarter struct {
	infra.BaseStarter
}

func (p *PropsStarter)Init(content infra.StarterContext) {
	// 初始化配置  在bootStarter实例化的时候就已经完成了，这里是验证有没有加载配置文件
	props = content.Props() //加到全局变量里面
	//fmt.Println("配置文件加载")
}

type SystemAccount struct {
	AccountNo string
	AccountName string
	UserId string
	Username string
}

var systemAccount *SystemAccount

var systemAccountOnce sync.Once

func GetSystemAccount() *SystemAccount {
	systemAccountOnce.Do(func() {
		systemAccount = new(SystemAccount)
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			panic(err)
		}
	})
	return systemAccount
}

func GetEnvelopeActivityLink() string {
	link := Props().GetDefault("envelope.link", "/v1/envelope/link")
	return link
}

func GetEnvelopeDomain() string {
	domain := Props().GetDefault("envelope.domain", "http://localhost")
	return domain
}