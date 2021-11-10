package services

const (
	DefaultBlessing = "恭喜发财，鸿富猪到！"
)

//订单类型：发布单、退款单
type OrderType int

const (
	OrderTypeSending OrderType = 1
	OrderTypeRefund  OrderType = 2
)

//支付状态：未支付，支付中，已支付，支付失败
//退款：未退款，退款中，已退款，退款失败
type PayStatus int

const (
	PayNothing PayStatus = 1
	Paying     PayStatus = 2
	Payed      PayStatus = 3
	PayFailure PayStatus = 4
)

//红包订单状态：创建、发布、过期、失效
type OrderStatus int

const (
	OrderCreate   OrderStatus = 1
	OrderSending  OrderStatus = 2
	OrderExpired  OrderStatus = 3
	OrderDisabled OrderStatus = 4
)

//红包类型：普通红包，碰运气红包
type EnvelopeType int

const (
	GeneralEnvelopeType = 1
	LuckyEnvelopeType   = 2
)
