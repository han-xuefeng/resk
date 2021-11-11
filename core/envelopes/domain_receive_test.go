package envelopes

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	//"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"study-gin/resk/services"
	_ "study-gin/resk/testx"
	"testing"
)

func TestRedEnvelopeService_Receive(t *testing.T) {
	accountService := services.GetAccountService()
	Convey("收红包测试用例", t, func() {
		accounts := make([]*services.AccountDTO, 0)
		size := 10
		for i := 0; i < size; i++ {
			account := services.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				Username:     "测试用户" + strconv.Itoa(i+1),
				Amount:       "2000",
				AccountName:  "测试账户" + strconv.Itoa(i+1),
				AccountType:  int(services.EnvelopeAccountType),
				CurrencyCode: "CNY",
			}
			//账户创建
			acDto, err := accountService.CreateAccount(account)
			So(err, ShouldBeNil)
			So(acDto, ShouldNotBeNil)
			accounts = append(accounts, acDto)
		}
		acDto := accounts[0]
		So(len(accounts), ShouldEqual, size)
		//2. 使用其中一个用户发送一个红包
		re := services.GetRedEnvelopeService()
		//发送普通红包
		goods := services.RedEnvelopeSendingDTO{
			UserId:       acDto.UserId,
			Username:     acDto.Username,
			EnvelopeType: services.GeneralEnvelopeType,
			Amount:       decimal.NewFromFloat(1.88),
			Quantity:     size,
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
		remainAmount := at.Amount
		//3. 使用发送红包数量的人收红包
		Convey("收普通红包", func() {
			for _, account := range accounts {
				rcv := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUserId:   account.UserId,
					RecvUsername: account.Username,
					AccountNo:    account.AccountNo,
				}
				item, err := re.Receive(rcv)
				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				So(item.Amount, ShouldEqual, at.AmountOne)
				remainAmount = remainAmount.Sub(at.AmountOne)
				So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())

			}
		})
		////收碰运气红包，作为作业留给同学们来实现
		goods.EnvelopeType = services.LuckyEnvelopeType
		goods.Amount = decimal.NewFromFloat(18.8)
		at, err = re.SendOut(goods)
		So(err, ShouldBeNil)
		So(at, ShouldNotBeNil)
		So(at.Link, ShouldNotBeEmpty)
		So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
		//验证每一个属性
		dto = at.RedEnvelopeGoodsDTO
		So(dto.Username, ShouldEqual, goods.Username)
		So(dto.UserId, ShouldEqual, goods.UserId)
		So(dto.Quantity, ShouldEqual, goods.Quantity)
		So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
		remainAmount = at.Amount
		re = services.GetRedEnvelopeService()
		Convey("收碰运气红包", func() {
			So(len(accounts), ShouldEqual, size)
			total := decimal.NewFromFloat(0)
			for i, account := range accounts {
				if i > 10 {
					break
				}
				rcv := services.RedEnvelopeReceiveDTO{
					EnvelopeNo:   at.EnvelopeNo,
					RecvUserId:   account.UserId,
					RecvUsername: account.Username,
					AccountNo:    account.AccountNo,
				}
				item, err := re.Receive(rcv)
				if item != nil {
					total = total.Add(item.Amount)
				}

				logrus.Info(i+1, " ", total.String(), " ", item.Amount.String())

				So(err, ShouldBeNil)
				So(item, ShouldNotBeNil)
				remainAmount = remainAmount.Sub(item.Amount)
				So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())

			}
			So(total.String(), ShouldEqual, goods.Amount.String())
		})

	})
}
