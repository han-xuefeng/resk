package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"study-gin/resk/services"
	_ "study-gin/resk/testx"
	"testing"
)

func TestRedEnvelopeService_SendOut(t *testing.T) {
	//发红包人的红包资金账户
	ac := services.GetAccountService()
	account := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "200",
		AccountName:  "测试账户",
		AccountType:  int(services.EnvelopeAccountType),
		CurrencyCode: "CNY",
	}
	re := services.GetRedEnvelopeService()
	Convey("准备资金账户", t, func() {
		//准备资金账户
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO, ShouldNotBeNil)
	})
	Convey("发红包", t, func() {
		Convey("发普通红包", func() {
			goods := services.RedEnvelopeSendingDTO{
				UserId:       account.UserId,
				Username:     account.Username,
				EnvelopeType: services.GeneralEnvelopeType,
				Amount:       decimal.NewFromFloat(8.88),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			at, err := re.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			//验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			q := decimal.NewFromFloat(float64(dto.Quantity))
			So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(q).String())
			//同学可以想一下，还需要验证哪些字段
		})
		Convey("发碰运气红包", func() {
			goods := services.RedEnvelopeSendingDTO{
				UserId:       account.UserId,
				Username:     account.Username,
				EnvelopeType: services.LuckyEnvelopeType,
				Amount:       decimal.NewFromFloat(88.8),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}

			at, err := re.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			//验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
			//同学可以想一下，还需要验证哪些字段
		})
	})
}
