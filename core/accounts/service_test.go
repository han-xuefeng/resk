package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"study-gin/resk/services"
	"testing"
)

func TestAccountService_CreateAccount(t *testing.T) {
	dto := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "100",
		AccountName:  "测试账户",
		AccountType:  2,
		CurrencyCode: "CNY",
	}
	service := new(accountService)
	Convey("账户创建", t, func() {
		rdto, err := service.CreateAccount(dto)
		So(err, ShouldBeNil)
		So(rdto, ShouldNotBeNil)
		So(rdto.Balance.String(), ShouldEqual, dto.Amount)
		So(rdto.UserId, ShouldEqual, dto.UserId)
		So(rdto.Username, ShouldEqual, dto.Username)
		So(rdto.Status, ShouldEqual, 1)
	})
}

//转账业务应用服务层测试用例
func TestAccountService_Transfer(t *testing.T) {

	Convey("转账", t, func() {
		//准备2个账户
		a1 := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户1",
			Amount:       "100",
			AccountName:  "测试账户1",
			AccountType:  2,
			CurrencyCode: "CNY",
		}
		a2 := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			Username:     "测试用户2",
			Amount:       "100",
			AccountName:  "测试账户2",
			AccountType:  2,
			CurrencyCode: "CNY",
		}
		service := new(accountService)
		adto1, err := service.CreateAccount(a1)
		So(err, ShouldBeNil)
		So(adto1, ShouldNotBeNil)
		adto2, err := service.CreateAccount(a2)
		So(err, ShouldBeNil)
		So(adto2, ShouldNotBeNil)

		////从账户1转入账户2一定金额，其中账户1的余额是足够的
		Convey("余额足够，从账户1转入账户2一定金额", func() {
			//转账过程需要2个交易的参与者：交易主体和交易对象
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
			amount := decimal.NewFromFloat(10)
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转出",
			}
			status, err := service.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)

			ra1 := service.GetAccount(adto1.AccountNo)
			So(ra1, ShouldNotBeNil)
			So(ra1.Balance.String(), ShouldEqual, adto1.Balance.Sub(amount).String())

		})
		//从账户1转入账户2一定金额，但余额不足，转账会失败
		Convey("余额不足，从账户1转入账户2一定金额", func() {
			//转账过程需要2个交易的参与者：交易主体和交易对象
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
			amount := adto1.Balance.Add(decimal.NewFromFloat(200))
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转出",
			}
			status, err := service.Transfer(dto)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferedStatusSufficientFunds)

			ra1 := service.GetAccount(adto1.AccountNo)
			So(ra1, ShouldNotBeNil)
			So(ra1.Balance.String(), ShouldEqual, adto1.Balance.String())
		})
		//给账户1储值
		Convey("给账户1储值", func() {
			//转账过程需要2个交易的参与者：交易主体和交易对象
			body := services.TradeParticipator{
				AccountNo: adto1.AccountNo,
				UserId:    adto1.UserId,
				Username:  adto1.Username,
			}
			target := body
			amount := decimal.NewFromFloat(10)
			dto := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				ChangeType:  services.AccountStoreValue,
				ChangeFlag:  services.FlagTransferIn,
				Decs:        "储值",
			}
			status, err := service.Transfer(dto)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferedStatusSuccess)

			ra1 := service.GetAccount(adto1.AccountNo)
			So(ra1, ShouldNotBeNil)
			So(ra1.Balance.String(), ShouldEqual, adto1.Balance.Add(amount).String())

		})
	})

}
