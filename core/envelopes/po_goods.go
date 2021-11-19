package envelopes

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"study-gin/resk/services"
	"time"
)

type RedEnvelopeGoods struct {
	Id             int64                `db:"id,omitempty"`         //自增ID
	EnvelopeNo     string               `db:"envelope_no,uni"`      //红包编号,红包唯一标识
	EnvelopeType   int                  `db:"envelope_type"`        //红包类型：普通红包，碰运气红包
	Username       sql.NullString       `db:"username"`             //用户名称
	UserId         string               `db:"user_id"`              //用户编号, 红包所属用户
	Blessing       sql.NullString       `db:"blessing"`             //祝福语
	Amount         decimal.Decimal      `db:"amount"`               //红包总金额
	AmountOne      decimal.Decimal      `db:"amount_one"`           //单个红包金额，碰运气红包无效
	Quantity       int                  `db:"quantity"`             //红包总数量
	RemainAmount   decimal.Decimal      `db:"remain_amount"`        //红包剩余金额额
	RemainQuantity int                  `db:"remain_quantity"`      //红包剩余数量
	ExpiredAt      time.Time            `db:"expired_at"`           //过期时间
	Status         services.OrderStatus `db:"status"`               //红包状态：0红包初始化，1启用，2失效
	OrderType      services.OrderType   `db:"order_type"`           //订单类型：发布单、退款单
	PayStatus      services.PayStatus   `db:"pay_status"`           //支付状态：未支付，支付中，已支付，支付失败
	CreatedAt      time.Time            `db:"created_at,omitempty"` //创建时间
	UpdatedAt      time.Time            `db:"updated_at,omitempty"` //更新时间
	OriginEnvelopeNo string               `db:"origin_envelope_no"`   //原关联订单号
}

func (po *RedEnvelopeGoods) ToDTO() *services.RedEnvelopeGoodsDTO {
	dto := &services.RedEnvelopeGoodsDTO{

		EnvelopeNo:     po.EnvelopeNo,
		EnvelopeType:   po.EnvelopeType,
		Username:       po.Username.String,
		UserId:         po.UserId,
		Blessing:       po.Blessing.String,
		Amount:         po.Amount,
		AmountOne:      po.AmountOne,
		Quantity:       po.Quantity,
		RemainAmount:   po.RemainAmount,
		RemainQuantity: po.RemainQuantity,
		ExpiredAt:      po.ExpiredAt,
		Status:         po.Status,
		OrderType:      po.OrderType,
		PayStatus:      po.PayStatus,
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
	}
	return dto
}

func (po *RedEnvelopeGoods) FromDTO(dto *services.RedEnvelopeGoodsDTO) {

	po.EnvelopeNo = dto.EnvelopeNo
	po.EnvelopeType = dto.EnvelopeType
	po.Username = sql.NullString{Valid: true, String: dto.Username}
	po.UserId = dto.UserId
	po.Blessing = sql.NullString{Valid: true, String: dto.Blessing}
	po.Amount = dto.Amount
	po.AmountOne = dto.AmountOne
	po.Quantity = dto.Quantity
	po.RemainAmount = dto.RemainAmount
	po.RemainQuantity = dto.RemainQuantity
	po.ExpiredAt = dto.ExpiredAt
	po.Status = dto.Status
	po.OrderType = dto.OrderType
	po.PayStatus = dto.PayStatus
	po.CreatedAt = dto.CreatedAt
	po.UpdatedAt = dto.UpdatedAt
}
