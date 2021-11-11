package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"study-gin/resk/services"
)

func main()  {
	c, err := rpc.Dial("tcp", ":8082")
	if err != nil {
		logrus.Panic(err)
	}
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(1),
		UserId:       "20iyiDYId85hl3VRurKK3yuQrAV",
		Username:     "测试用户",
		EnvelopeType: services.GeneralEnvelopeType,
		Quantity:     2,
		Blessing:     "",
	}
	out := &services.RedEnvelopeActivity{}
	err = c.Call("EnvelopeRpc.SendOut", in, out)
	if err != nil {
		logrus.Panic(err)
	}
	c.Close()
	logrus.Info(out)
}
