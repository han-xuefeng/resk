package envelopes

import (
	"context"
	"github.com/tietang/dbx"
	"path"
	"study-gin/resk/core/accounts"
	"study-gin/resk/infra/base"
	"study-gin/resk/services"
)

// 发红包业务领域代码
func (d *goodsDomain) SendOut(goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包对象
	d.Create(goods)
	// 创建活动
	activity = new(services.RedEnvelopeActivity)
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)
	accountDomain := accounts.NewAccountDomain()
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		//保存红包商品
		id, err := d.Save(ctx)
		if id <= 0 || err != nil {
			return err
		}
		body := services.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			Username:  systemAccount.Username,
			UserId:    systemAccount.UserId,
		}

		transfer := services.AccountTransferDTO{
			TradeBody:   body,
			TradeTarget: target,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferOut,
			Decs:        "红包金额支付",
		}
		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status == services.TransferedStatusSuccess {
			return nil
		}
		//3. 将扣减的红包总金额转入红包中间商的红包资金账户
		//入账
		transfer = services.AccountTransferDTO{
			TradeBody:   target,
			TradeTarget: body,
			TradeNo:     d.EnvelopeNo,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeIncoming,
			ChangeFlag:  services.FlagTransferIn,
			Decs:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if status == services.TransferedStatusSuccess {
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	//扣减金额没有问题，返回活动

	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()
	return
}