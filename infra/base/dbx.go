package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
	"study-gin/resk/infra"
)

var database *dbx.Database

func DbxDatabase() *dbx.Database {
	return database
}


type DbxDatabaseStarter struct {
	infra.BaseStarter
}

func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := Props()
	settings := dbx.Settings{}
		kvs.Unmarshal(conf, &settings, "mysql")

	dbx, err := dbx.Open(settings)

	if err != nil {
		panic(nil)
	}
	logrus.Info(dbx.Ping())
	database = dbx
}