package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"study-gin/resk/services"
	"testing"
)

func TestAccountDomain_Create(t *testing.T) {
	dto := services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户",
		Balance:  decimal.NewFromFloat(0),
		Status:   1,
	}
	domain := new(accountDomain)
	Convey("账户创建", t, func() {
		rdto, err := domain.Create(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Balance.String())
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.Status, ShouldEqual, dto.Status)
	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	//2个账户，交易主体账户要有余额
	adto1 := &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户1",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	adto2 := &services.AccountDTO{
		UserId:   ksuid.New().Next().String(),
		Username: "测试用户2",
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	domain := accountDomain{}
	Convey("转账测试", t, func() {
		//账户1的准备
		dto1, err := domain.Create(*adto1)
		So(err, ShouldBeNil)
		So(dto1, ShouldNotBeNil)
		So(dto1.Balance.String(), ShouldEqual, adto1.Balance.String())
		So(dto1.UserId, ShouldEqual, adto1.UserId)
		So(dto1.Username, ShouldEqual, adto1.Username)
		So(dto1.Status, ShouldEqual, adto1.Status)
		adto1 = dto1
		//账户2的准备
		dto2, err := domain.Create(*adto2)
		So(err, ShouldBeNil)
		So(dto2, ShouldNotBeNil)
		So(dto2.Balance.String(), ShouldEqual, adto2.Balance.String())
		So(dto2.UserId, ShouldEqual, adto2.UserId)
		So(dto2.Username, ShouldEqual, adto2.Username)
		So(dto2.Status, ShouldEqual, adto2.Status)
		adto2 = dto2

		Convey("余额充足，金额转入其他账户", func() {
			amount := decimal.NewFromFloat(1)
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转账",
			}
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)
			//实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(),
				ShouldEqual,
				adto1.Balance.Sub(amount).String())

		})
		//2. 余额不足，金额转出
		Convey("余额不足，金额转出", func() {
			amount := adto1.Balance
			amount = amount.Add(decimal.NewFromFloat(200))
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转账",
			}
			status, err := domain.Transfer(dto)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferedStatusSufficientFunds)
			//实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(),
				ShouldEqual,
				adto1.Balance.String())

		})
		//3. 充值
		Convey("充值", func() {
			amount := decimal.NewFromFloat(11.1)
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			target := services.TradeParticipator{
				AccountNo: adto2.AccountNo,
				UserId:    adto2.UserId,
				Username:  adto2.Username,
			}
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				Amount:      amount,
				ChangeType:  services.AccountStoreValue,
				ChangeFlag:  services.FlagTransferIn,
				Decs:        "储值",
			}
			status, err := domain.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)
			//实际余额更新后的预期值
			a2 := domain.GetAccount(adto1.AccountNo)
			So(a2, ShouldNotBeNil)
			So(a2.Balance.String(),
				ShouldEqual,
				adto1.Balance.Add(amount).String())

		})
	})
}
