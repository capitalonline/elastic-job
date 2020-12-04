package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mongodb-job/config"
	"log"
	"time"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

type (
	// TODO 用于填充系统默认数据的接口
	Seeder interface {
		Seed() error
	}
	// TODO 模型接口
	Model interface {
		Store() error
		Update() error
		ToString() (string, error)
	}
)

var Engine *xorm.Engine
var xormerr error


func Connection(){
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		config.Conf.Database.User,
		config.Conf.Database.Pass,
		config.Conf.Database.Host,
		config.Conf.Database.Port,
		config.Conf.Database.Name,
		config.Conf.Database.Char,
	)
	Engine, xormerr = xorm.NewEngine("mysql", dsn)
	if xormerr != nil {
		log.Fatal("dial database connections failed")
	}
	if Engine != nil {
		Engine.SetMaxIdleConns(30)
		Engine.SetMaxOpenConns(50)
		//Engine.ShowSQL(true)
		Engine.Logger().SetLevel(xlog.LOG_DEBUG)
		Engine.SetConnMaxLifetime(time.Second * 30)
	}

	if xormerr = Engine.Ping(); xormerr != nil {
		log.Fatal("test database connection failed")
	}
}