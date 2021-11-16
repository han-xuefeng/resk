package jobs

import (
	"github.com/go-redsync/redsync"
	//"github.com/gomodule/redigo/redis"
	log "github.com/sirupsen/logrus"
	"study-gin/resk/infra"
	"time"
)

type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
	mutex  *redsync.Mutex
}

func (r *RefundExpiredJobStarter)Init(ctx infra.StarterContext)  {
	d := ctx.Props().GetDurationDefault("jobs.refund.interval", time.Minute)
	r.ticker = time.NewTicker(d)
	//maxIdle := ctx.Props().GetIntDefault("redis.maxIdle", 2)
	//maxActive := ctx.Props().GetIntDefault("redis.maxActive", 5)
	//timeout := ctx.Props().GetDurationDefault("redis.timeout", 20*time.Second)
	//addr := ctx.Props().GetDefault("redis.addr", "127.0.0.1:6379")
	//pools := make([]redsync.Pool, 0)
	//pool := &redis.Pool{
	//	MaxIdle:     maxIdle,
	//	MaxActive:   maxActive,
	//	IdleTimeout: timeout,
	//	Dial: func() (conn redis.Conn, e error) {
	//		return redis.Dial("tcp", addr)
	//	},
	//}
	//pools = append(pools, pool)
	//rsync := redsync.New(pools)
}

func (r *RefundExpiredJobStarter)Start(ctx infra.StarterContext)  {
	go func() {
		for {
			c := r.ticker.C
			log.Debug("过期红包退款开始。。。", c)
			// 红包逻辑退款的代码
		}
	}()
}

func (r *RefundExpiredJobStarter)Stop(ctx infra.StarterContext)  {
	r.ticker.Stop()
}