package resk

import (
	"study-gin/resk/apis/gorpc"
	_ "study-gin/resk/apis/gorpc"
	_ "study-gin/resk/apis/web"
	_ "study-gin/resk/core/accounts"
	_ "study-gin/resk/core/envelopes"
	"study-gin/resk/infra"
	"study-gin/resk/infra/base"
	"study-gin/resk/jobs"
)

//这里面注册
func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRPCStarter{})
	infra.Register(&gorpc.GoRpcApiStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
}
