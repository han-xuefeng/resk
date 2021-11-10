package envelopes

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct {
	runner *dbx.TxRunner
}

func (dao *RedEnvelopeItemDao) GetOne(itemNo string) *RedEnvelopeItem {
	po := &RedEnvelopeItem{ItemNo: itemNo}
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

//红包订单详情数据的写入 Insert

func (dao *RedEnvelopeItemDao) Insert(form *RedEnvelopeItem) (int64, error) {
	rs, err := dao.runner.Insert(form)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *RedEnvelopeItemDao) FindItems(envelopeNo string) []*RedEnvelopeItem {
	items := make([]*RedEnvelopeItem, 0)
	sql := "select * from red_envelope_item where envelope_no=?"
	err := dao.runner.Find(items, sql, envelopeNo)
	if err != nil {
		logrus.Error(err)
		return items
	}
	return items
}