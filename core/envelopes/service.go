package envelopes

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"study-gin/resk/infra/base"
	"study-gin/resk/services"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.IRedEnvelopeService = new(redEnvelopeService)
	})
}

type redEnvelopeService struct {
}
//收红包
func (r *redEnvelopeService) Receive(dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	//参数效验
	if err = base.ValidateStruct(&dto); err != nil {
		return nil, err
	}
	//获取当前收红包用户的账户信息
	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.RecvUserId)
	if account == nil {
		return nil, errors.New("红包资金账户不存在：user_id=" + dto.RecvUserId)
	}
	dto.AccountNo = account.AccountNo
	//进行尝试收红包
	domain := goodsDomain{}
	itemDomain := itemDomain{}
	item = itemDomain.GetByUser(dto.RecvUserId, dto.EnvelopeNo)
	if item != nil {
		return item, nil
	}
	item, err = domain.Receive(context.Background(), dto)
	return item, err
}
func (r *redEnvelopeService) Refund(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	panic("implement me")
}

func (r *redEnvelopeService) Get(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	panic("implement me")
}

//发红包
func (r *redEnvelopeService) SendOut(
	dto services.RedEnvelopeSendingDTO) (activity *services.RedEnvelopeActivity, err error) {
	//验证
	if err = base.ValidateStruct(&dto); err != nil {
		return activity, err
	}
	//获取红包发送人的资金账户信息
	account := services.GetAccountService().GetEnvelopeAccountByUserId(dto.UserId)
	if account == nil {
		return nil, errors.New("用户账户不存在：" + dto.UserId)
	}
	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo

	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}
	//执行发送红包的逻辑
	domain := new(goodsDomain)
	activity, err = domain.SendOut(*goods)
	if err != nil {
		log.Error(err)
	}

	return activity, err
}

