package accounts

import (
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"study-gin/resk/infra/base"
	"study-gin/resk/services"
)

// 有状态的, 每次使用都要实例化
type accountDomain struct {
	account Account
	accountLog AccountLog
}

// createAccountLogNo 创建logNo的逻辑
func (domain *accountDomain) createAccountLogNo() {
	domain.accountLog.LogNo = ksuid.New().Next().String()
}

// createAccountNo 生成accountNo的逻辑
func (domain *accountDomain) createAccountNo() {
	domain.account.AccountNo = ksuid.New().Next().String()
}

// 创建账户流水的记录
func (domain *accountDomain) createAccountLog() {
	domain.accountLog = AccountLog{}
	domain.createAccountLogNo()
	domain.accountLog.TradeNo = domain.accountLog.LogNo

	domain.accountLog.AccountNo = domain.account.AccountNo
	domain.accountLog.UserId = domain.account.UserId
	domain.accountLog.Username = domain.account.Username.String

	domain.accountLog.TargetAccountNo = domain.account.AccountNo
	domain.accountLog.TargetUserId = domain.account.UserId
	domain.accountLog.TargetUsername = domain.account.Username.String

	domain.accountLog.Amount = domain.account.Balance
	domain.accountLog.Balance = domain.account.Balance

	domain.accountLog.Decs = "账户创建"
	domain.accountLog.ChangeType = services.AccountCreated
	domain.accountLog.ChangeFlag = services.FlagAccountCreated
}

//账户创建的业务逻辑
func (domain *accountDomain) Create(
	dto services.AccountDTO) (*services.AccountDTO, error) {
	//创建账户持久化对象
	domain.account = Account{}
	domain.account.FromDTO(&dto)
	domain.createAccountNo()
	domain.account.Username.Valid = true
	//创建账户流水持久化对象
	domain.createAccountLog()
	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}
	var rdto *services.AccountDTO
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner
		//插入账户数据
		id, err := accountDao.Insert(&domain.account)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户失败")
		}
		//如果插入成功，就插入流水数据
		id, err = accountLogDao.Insert(&domain.accountLog)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户流水失败")
		}
		domain.account = *accountDao.GetOne(domain.account.AccountNo)
		return nil
	})
	rdto = domain.account.ToDTO()
	return rdto, err

}

func (domain *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
	//如果交易变化是支出，修正amount
	amount := dto.Amount
	if dto.ChangeFlag == services.FlagTransferOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}

	//创建账户流水记录
	domain.accountLog = AccountLog{}
	domain.accountLog.FromTransferDTO(&dto)
	domain.createAccountLogNo()
	//检查余额是否足够和更新余额：通过乐观锁来验证，更新余额的同时来验证余额是否足够
	//更新成功后，写入流水记录
	err = base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		accountLogDao := AccountLogDao{runner: runner}

		rows, err := accountDao.UpdateBalance(dto.TradeBody.AccountNo, amount)
		if err != nil {
			status = services.TransferedStatusFailure
			return err
		}
		if rows <= 0 && dto.ChangeFlag == services.FlagTransferOut {
			status = services.TransferedStatusSufficientFunds
			return errors.New("余额不足")
		}
		account := accountDao.GetOne(dto.TradeBody.AccountNo)
		if account == nil {
			return errors.New("账户出错")
		}
		domain.account = *account
		domain.accountLog.Balance = domain.account.Balance
		id, err := accountLogDao.Insert(&domain.accountLog)
		if err != nil || id <= 0 {
			status = services.TransferedStatusFailure
			return errors.New("账户流水创建失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
	} else {
		status = services.TransferedStatusSuccess
	}

	return status, err
}

//根据账户编号来查询账户信息
func (domain *accountDomain) GetAccount(accountNo string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetOne(accountNo)
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

//根据用户ID来查询红包账户信息
func (domain *accountDomain) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetByUserId(userId, int(services.EnvelopeAccountType))
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()

}

func (domain *accountDomain) GetAccountByUserIdAndType(userId string, accountType services.AccountType) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetByUserId(userId, int(accountType))
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

//根据流水ID来查询账户流水
func (domain *accountDomain) GetAccountLog(logNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetOne(logNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}

//根据交易编号来查询账户流水
func (domain *accountDomain) GetAccountLogByTradeNo(tradeNo string) *services.AccountLogDTO {
	dao := AccountLogDao{}
	var log *AccountLog
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao.runner = runner
		log = dao.GetByTradeNo(tradeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if log == nil {
		return nil
	}
	return log.ToDTO()
}
