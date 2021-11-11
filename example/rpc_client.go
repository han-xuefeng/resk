package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"study-gin/resk/services"
)

func main() {
	c, err := rpc.Dial("tcp", ":8082")
	if err != nil {
		logrus.Panic(err)
	}
	//sendout(c)
	receive(c)

}

func receive(c *rpc.Client) {
	in := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   "20mLEDj74jnw9QSzX3AxmQXeDZo",
		RecvUserId:   "20mLENrRs3cnsS78mgPoaxI4BZ8",
		RecvUsername: "测试用户10",
		AccountNo:    "",
	}
	out := &services.RedEnvelopeItemDTO{}
	err := c.Call("EnvelopeRpc.Receive", in, out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("%+v", out)
}
func sendout(c *rpc.Client) {
	in := services.RedEnvelopeSendingDTO{
		Amount:       decimal.NewFromFloat(1),
		UserId:       "20mLENrRs3cnsS78mgPoaxI4BZ8",
		Username:     "测试用户",
		EnvelopeType: services.GeneralEnvelopeType,
		Quantity:     2,
		Blessing:     "",
	}
	out := &services.RedEnvelopeActivity{}
	err := c.Call("EnvelopeRpc.SendOut", in, &out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("%+v", out)
}
