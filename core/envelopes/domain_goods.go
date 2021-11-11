package envelopes

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"study-gin/resk/infra/base"
	"study-gin/resk/services"
	"time"
)

// 商品领域层
type goodsDomain struct {
	RedEnvelopeGoods
	item itemDomain
}
// 生成一个红包编号
func (d *goodsDomain) createEnvelopeNo() {
	d.EnvelopeNo = ksuid.New().Next().String()
}

// 创建一个红包商品对象
func (d *goodsDomain) Create(goods services.RedEnvelopeGoodsDTO) {
	d.RedEnvelopeGoods.FromDTO(&goods)
	d.RemainQuantity = goods.Quantity
	d.Username.Valid = true
	d.Blessing.Valid = true
	if d.EnvelopeType == services.GeneralEnvelopeType {
		d.Amount = goods.AmountOne.Mul(decimal.NewFromFloat(float64(goods.Quantity)))
	}
	if d.EnvelopeType == services.LuckyEnvelopeType {
		d.AmountOne = decimal.NewFromFloat(0)
	}
	d.RemainAmount = d.Amount
	//过期时间
	d.ExpiredAt = time.Now().Add(24 * time.Hour)
	d.Status = services.OrderCreate
	d.createEnvelopeNo()
}
// 保存
func (d *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		id, err = dao.Insert(&d.RedEnvelopeGoods)
		return err
	})
	return id, err
}

// 创建对象并且保存
func (d *goodsDomain) CreateAndSave(ctx context.Context, goods services.RedEnvelopeGoodsDTO) (id int64, err error) {
	d.Create(goods)
	return d.Save(ctx)
}

//查询红包商品信息
func (d *goodsDomain) Get(envelopeNo string) (goods *RedEnvelopeGoods) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
	return goods
}