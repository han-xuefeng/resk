package envelopes

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"study-gin/resk/core/accounts"
	"study-gin/resk/infra/algo"
	"study-gin/resk/infra/base"
	"study-gin/resk/services"
)
var multiple = decimal.NewFromFloat(100.0)
// 收红包业务
func (d *goodsDomain) Receive(
	ctx context.Context,
	dto services.RedEnvelopeReceiveDTO)(item *services.RedEnvelopeItemDTO, err error){
	// 1.创建收红包的订单明细
	d.preCreateItem(dto)
	// 2. 查询出当前红包剩余数量和剩余金额
	goods := d.Get(dto.EnvelopeNo)
	//3. 效验剩余红包和剩余金额：
	//- 如果没有剩余，直接返回无可用红包金额
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		log.Errorf("%+v", goods)
		return nil, errors.New("没有足够的红包和金额了")
	}
	//4. 使用红包算法计算红包金额
	nextAmount := d.nextAmount(goods)
	err = base.Tx(func(runner *dbx.TxRunner) error {
		//5. 使用乐观锁更新语句，尝试更新剩余数量和剩余金额：
		dao := RedEnvelopeGoodsDao{
			runner: runner,
		}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		if err != nil || rows <= 0 {
			log.Errorf("rows=%d,%s", rows, err.Error())
			return errors.New("没有足够的红包和金额")
		}
		//6. 保存订单明细数据
		d.item.Quantity = 1
		d.item.PayStatus = int(services.Paying)
		d.item.AccountNo = dto.AccountNo
		d.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		d.item.Amount = nextAmount
		desc := goods.Username.String + "的" + services.EnvelopeTypes[services.EnvelopeType(goods.EnvelopeType)]
		d.item.Desc = desc
		txCtx := base.WithValueContext(ctx, runner)
		_, err = d.item.Save(txCtx)
		if err != nil {
			log.Error(err)
			return err
		}
		//7. 将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
		// transfer
		status, err := d.transfer(txCtx, dto)
		if status == services.TransferedStatusSuccess {
			return nil
		}
		return err
	})
	return d.item.ToDTO(), err
}

func (d *goodsDomain) transfer(
	ctx context.Context,
	dto services.RedEnvelopeReceiveDTO) (status services.TransferedStatus, err error) {
	systemAccount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		Username:  systemAccount.Username,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		Username:  dto.RecvUsername,
	}
	transfer := services.AccountTransferDTO{
		TradeBody:   body,
		TradeTarget: target,
		TradeNo:     dto.EnvelopeNo,
		Amount:      d.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包收入",
	}
	adomain := accounts.NewAccountDomain()
	return adomain.TransferWithContextTx(ctx, transfer)
}


// 创建收红包订单明细
func (d *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	d.item.AccountNo = dto.AccountNo
	d.item.EnvelopeNo = dto.EnvelopeNo
	d.item.RecvUsername = sql.NullString{String: dto.RecvUsername, Valid:true}
	d.item.RecvUserId = dto.RecvUserId
	d.item.createItemNo()
}

// 计算下一个红包的金额
func (d *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		return goods.AmountOne
	} else if goods.EnvelopeType == services.LuckyEnvelopeType {
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	} else {
		log.Error("不支持的红包类型")
	}
	return amount

}