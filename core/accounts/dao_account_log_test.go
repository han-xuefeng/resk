package accounts

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"study-gin/resk/infra/base"
	"study-gin/resk/services"
	"testing"
	_ "study-gin/resk/testx"
)

func TestAccountLogDao(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountLogDao{
			runner: runner,
		}
		Convey("通过Log编号查询账户流水数据", t, func() {
			a := &AccountLog{
				LogNo:      ksuid.New().Next().String(),
				TradeNo:    ksuid.New().Next().String(),
				Status:     1,
				AccountNo:  ksuid.New().Next().String(),
				UserId:     ksuid.New().Next().String(),
				Username:   "测试用户",
				Amount:     decimal.NewFromFloat(1),
				Balance:    decimal.NewFromFloat(100),
				ChangeFlag: services.FlagAccountCreated,
				ChangeType: services.AccountCreated,
			}
			//通过log no来查询
			Convey("通过log no来查询", func() {
				id, err := dao.Insert(a)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)
				na := dao.GetOne(a.LogNo)
				So(na, ShouldNotBeNil)
				So(na.Balance.String(), ShouldEqual, a.Balance.String())
				So(na.Amount.String(), ShouldEqual, a.Amount.String())
				So(na.CreatedAt, ShouldNotBeNil)
			})

			//通过trade no来查询
			Convey("通过trade no来查询", func() {
				id, err := dao.Insert(a)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)
				na := dao.GetByTradeNo(a.TradeNo)
				So(na, ShouldNotBeNil)
				So(na.Balance.String(), ShouldEqual, a.Balance.String())
				So(na.Amount.String(), ShouldEqual, a.Amount.String())
				So(na.CreatedAt, ShouldNotBeNil)
			})
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}