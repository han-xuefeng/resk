package web

import (
	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
	"study-gin/resk/infra/base"
	"study-gin/resk/services"
)

type AccountApi struct {
	service services.AccountService
}

func (a *AccountApi) Init() {
	groupRouter := base.Iris().Party("/v1/account")
	groupRouter.Post("create", createHandler)
}

func createHandler(ctx iris.Context) {
	account := services.AccountCreatedDTO{}
	err := ctx.ReadJSON(&account)
	r := base.Res{
		Code: base.ResCodeOk,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamsError
		r.Message = err.Error()
		ctx.JSON(&r)
		logrus.Error(err)
		return
	}
	//执行创建账户的代码
	service := services.GetAccountService()
	dto, err := service.CreateAccount(account)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		logrus.Error(err)
	}
	r.Data = dto
	ctx.JSON(&r)
}