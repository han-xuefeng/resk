package accounts

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDao struct {
	runner *dbx.TxRunner
}

// GetOne 查询数据库持久化对象
func (dao *AccountDao) GetOne(accountNo string) *Account {
	a := &Account{AccountNo: accountNo}
	ok,err := dao.runner.GetOne(a)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	fmt.Println(a)
	return a
}

// GetByUserId 通过用户ID和账户类型来查询账户信息
func (dao *AccountDao) GetByUserId(
	userId string,
	accountType int) *Account {
	a := &Account{}
	sql := "select * from account" +
		" where user_id=? and account_type=? "
	ok, err := dao.runner.Get(a, sql, userId, accountType)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

// Insert 账户数据的插入
func (dao *AccountDao) Insert(a *Account) (id int64, err error) {
	rs, err := dao.runner.Insert(a)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

// UpdateBalance 账户余额的更新
//amount 如果是负数，就是扣减；如果是正数，就是增加
func (dao *AccountDao) UpdateBalance(
	accountNo string,
	amount decimal.Decimal) (rows int64, err error) {
	sql := "update account " +
		" set balance=balance+CAST(? AS DECIMAL(30,6))" +
		" where account_no=? " +
		" and balance>=-1*CAST(? AS DECIMAL(30,6)) "
	rs, err := dao.runner.Exec(sql,
		amount.String(),
		accountNo,
		amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

// UpdateStatus 账户状态更新
func (dao *AccountDao) UpdateStatus(
	accountNo string,
	status int) (rows int64, err error) {
	sql := "update account set status=? " +
		" where account_no=? "
	rs, err := dao.runner.Exec(sql, status, accountNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()

}