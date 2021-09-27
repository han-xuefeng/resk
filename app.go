package resk

import (
	"study-gin/resk/infra"
	"study-gin/resk/infra/base"
)

//这里面注册
func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.IrisServerStarter{})
}