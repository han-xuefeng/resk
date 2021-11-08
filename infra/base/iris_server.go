package base

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	irisrecover "github.com/kataras/iris/v12/middleware/recover"
	"github.com/sirupsen/logrus"
	"study-gin/resk/infra"
	"time"
)

var irisApplication *iris.Application

func Iris() *iris.Application {
	Check(irisApplication)
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (i *IrisServerStarter) Init(ctx infra.StarterContext) {
	fmt.Println("iris Init()")
	irisApplication = initIris()
	logger := irisApplication.Logger()
	logger.Install(logrus.StandardLogger())
}

func (i *IrisServerStarter) Setup(ctx infra.StarterContext) {
	fmt.Println("iris SetUp()")
}
func (i *IrisServerStarter) Start(ctx infra.StarterContext) {
	fmt.Println("iris Start()")
	//和logrus日志级别保持一致
	Iris().Logger().SetLevel(ctx.Props().GetDefault("log.level", "info"))

	//把路由打印到控制台
	routers := Iris().GetRoutes()
	fmt.Println("打印iris的路由信息")
	for _, router := range routers {
		logrus.Info(router.Trace)
	}

	port := Props().GetDefault("app.server.port", "18080")
	fmt.Println("iris要启动了")
	Iris().Run(iris.Addr(":" + port))
}

func (i *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application{
	//iris
	app := iris.New()
	app.Use(irisrecover.New())
	cfg := logger.Config{
		Status: true,
		IP:     true,
		Method: true,
		Path:   true,
		Query:  true,
		LogFunc: func(now time.Time, latency time.Duration,
			status, ip, method, path string,
			message interface{},
			headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s",
				now.Format("2006-01-02.15:04:05.000000"),
				latency.String(), status, ip, method, path, headerMessage, message,
			)
		},
	}
	app.Use(logger.New(cfg))
	return app
}