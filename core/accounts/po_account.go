package accounts

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"study-gin/resk/services"
	"time"
)
type Account struct {
	Id           int64           `db:"id,omitempty"`         //账户ID
	AccountNo    string          `db:"account_no,uni"`       //账户编号,账户唯一标识
	AccountName  string          `db:"account_name"`         //账户名称,用来说明账户的简短描述,账户对应的名称或者命名，比如xxx积分、xxx零钱
	AccountType  int             `db:"account_type"`         //账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户
	CurrencyCode string          `db:"currency_code"`        //货币类型编码：CNY人民币，EUR欧元，USD美元 。。。
	UserId       string          `db:"user_id"`              //用户编号, 账户所属用户
	Username     sql.NullString  `db:"username"`             //用户名称
	Balance      decimal.Decimal `db:"balance"`              //账户可用余额
	Status       int             `db:"status"`               //账户状态，账户状态：0账户初始化，1启用，2停用
	CreatedAt    time.Time       `db:"created_at,omitempty"` //创建时间
	UpdatedAt    time.Time       `db:"updated_at,omitempty"` //更新时间
}

// ToDTO ,omitempty
func (po *Account) ToDTO() *services.AccountDTO {
	dto := &services.AccountDTO{}
	dto.AccountNo = po.AccountNo
	dto.AccountName = po.AccountName
	dto.AccountType = po.AccountType
	dto.CurrencyCode = po.CurrencyCode
	dto.UserId = po.UserId
	dto.Username = po.Username.String
	dto.Balance = po.Balance
	dto.Status = po.Status
	dto.CreatedAt = po.CreatedAt
	dto.UpdatedAt = po.UpdatedAt
	return dto
}

func (po *Account) FromDTO(dto *services.AccountDTO) {
	po.AccountNo = dto.AccountNo
	po.AccountName = dto.AccountName
	po.AccountType = dto.AccountType
	po.CurrencyCode = dto.CurrencyCode
	po.UserId = dto.UserId
	po.Username = sql.NullString{Valid: true, String: dto.Username}
	po.Balance = dto.Balance
	po.Status = dto.Status
	po.CreatedAt = dto.CreatedAt
	po.UpdatedAt = dto.UpdatedAt
}
