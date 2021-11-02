package services

// TransferedStatus 转账状态
type TransferedStatus int8

// TransferedStatusFailure 转账失败
const TransferedStatusFailure TransferedStatus = -1

// TransferedStatusSufficientFunds 余额不足
const TransferedStatusSufficientFunds TransferedStatus = 0

// TransferedStatusSuccess 转账成功
const TransferedStatusSuccess TransferedStatus = 1

// ChangeType 转账的类型：0=创建账户 >=1进账 <=- 支出
type ChangeType int8

const (
	// AccountCreated 账户创建
	AccountCreated ChangeType = 0
	// AccountStoreValue 储值
	AccountStoreValue ChangeType = 1
	// EnvelopeOutgoing 红包资金的支出
	EnvelopeOutgoing ChangeType = -2
	// EnvelopeIncoming 红包资金的收入
	EnvelopeIncoming ChangeType = 2
	// EnvelopExpiredRefund 红包过期退款
	EnvelopExpiredRefund ChangeType = 3
)

// ChangeFlag 资金交易变化标识
type ChangeFlag int8

const (
	// FlagAccountCreated 创建账户=0
	FlagAccountCreated ChangeFlag = 0
	// FlagTransferOut 支出=-1
	FlagTransferOut ChangeFlag = -1
	// FlagTransferIn 收入=1
	FlagTransferIn ChangeFlag = 1
)

//账户类型
type AccountType int8

const (
	EnvelopeAccountType       AccountType = 1
	SystemEnvelopeAccountType AccountType = 2
)

