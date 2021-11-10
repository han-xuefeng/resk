package accounts

import (
	"context"
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

func NewAccountDomain() * accountDomain {
	return &accountDomain{}
}

// createAccountLogNo 创建logNo的逻辑
func (a *accountDomain) createAccountLogNo() {
	a.accountLog.LogNo = ksuid.New().Next().String()
}

// createAccountNo 生成accountNo的逻辑
func (a *accountDomain) createAccountNo() {
	a.account.AccountNo = ksuid.New().Next().String()
}

// 创建账户流水的记录
func (a *accountDomain) createAccountLog() {
	a.accountLog = AccountLog{}
	a.createAccountLogNo()
	a.accountLog.TradeNo = a.accountLog.LogNo

	a.accountLog.AccountNo = a.account.AccountNo
	a.accountLog.UserId = a.account.UserId
	a.accountLog.Username = a.account.Username.String

	a.accountLog.TargetAccountNo = a.account.AccountNo
	a.accountLog.TargetUserId = a.account.UserId
	a.accountLog.TargetUsername = a.account.Username.String

	a.accountLog.Amount = a.account.Balance
	a.accountLog.Balance = a.account.Balance

	a.accountLog.Decs = "账户创建"
	a.accountLog.ChangeType = services.AccountCreated
	a.accountLog.ChangeFlag = services.FlagAccountCreated
}

//账户创建的业务逻辑
func (a *accountDomain) Create(
	dto services.AccountDTO) (*services.AccountDTO, error) {
	//创建账户持久化对象
	a.account = Account{}
	a.account.FromDTO(&dto)
	a.createAccountNo()
	a.account.Username.Valid = true
	//创建账户流水持久化对象
	a.createAccountLog()
	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}
	var rdto *services.AccountDTO
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		accountLogDao.runner = runner
		//插入账户数据
		id, err := accountDao.Insert(&a.account)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户失败")
		}
		//如果插入成功，就插入流水数据
		id, err = accountLogDao.Insert(&a.accountLog)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户流水失败")
		}
		a.account = *accountDao.GetOne(a.account.AccountNo)
		return nil
	})
	rdto = a.account.ToDTO()
	return rdto, err

}

func (a *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		status, err = a.TransferWithContextTx(ctx, dto)
		return err
	})
	return status, err
}

//必须在base.TX事务块里面运行，不能单独运行
func (a *accountDomain) TransferWithContextTx(ctx context.Context, dto services.AccountTransferDTO) (status services.TransferedStatus, err error) {
	//如果交易变化是支出，修正amount
	amount := dto.Amount
	if dto.ChangeFlag == services.FlagTransferOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}

	//创建账户流水记录
	a.accountLog = AccountLog{}
	a.accountLog.FromTransferDTO(&dto)
	a.createAccountLogNo()
	//检查余额是否足够和更新余额：通过乐观锁来验证，更新余额的同时来验证余额是否足够
	//更新成功后，写入流水记录
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
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
		a.account = *account
		a.accountLog.Balance = a.account.Balance
		id, err := accountLogDao.Insert(&a.accountLog)
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
func (a *accountDomain) GetAccount(accountNo string) *services.AccountDTO {
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
func (a *accountDomain) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
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

func (a *accountDomain) GetAccountByUserIdAndType(userId string, accountType services.AccountType) *services.AccountDTO {
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
func (a *accountDomain) GetAccountLog(logNo string) *services.AccountLogDTO {
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
func (a *accountDomain) GetAccountLogByTradeNo(tradeNo string) *services.AccountLogDTO {
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
