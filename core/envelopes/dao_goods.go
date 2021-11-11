package envelopes

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"study-gin/resk/services"
	"time"
)
// 查询 根据红包编号
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	po := &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(po)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return po
}

type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

// 插入
func (dao *RedEnvelopeGoodsDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	rs ,err := dao.runner.Insert(po)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}
// 更新余额和数量
func (dao *RedEnvelopeGoodsDao) UpdateBalance(
	envelopeNo string, amount decimal.Decimal) (int64, error) {
	sql := "update red_envelope_goods " +
		" set remain_amount=remain_amount-CAST(? AS DECIMAL(30,6)), " +
		" remain_quantity=remain_quantity-1 " +
		" where envelope_no=? " +
		//最重要的，乐观锁的关键
		" and remain_quantity>0" +
		" and remain_amount >= CAST(? AS DECIMAL(30,6)) "

	rs, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()

}

// 更新订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(
	envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := " update red_envelope_goods" +
		" set order_status=? " +
		" where envelope_no=?"

	rs, err := dao.runner.Exec(sql, status, envelopeNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()

}
// 过期 把过期的所有红包都查询出来
func (dao *RedEnvelopeGoodsDao) FindExpired(
	offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := " select * from red_envelope_goods " +
		" where expired_at=? " +
		" limit ?,?"
	err := dao.runner.Find(&goods, sql, now, offset, size)
	if err != nil {
		logrus.Error(err)
	}
	return goods
}