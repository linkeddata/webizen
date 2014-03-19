package webizen

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
)

type User struct {
	Id          int64
	Uri         string    `xorm:"unique varbinary(500) not null"`
	CreatedTime time.Time `xorm:"index created"`
	UpdatedTime time.Time `xorm:"index updated"`
}

type UserName struct {
	User int64  `xorm:"index"`
	Name string `xorm:"index varchar(1000) not null"`
}

type UserImage struct {
	User  int64  `xorm:"index"`
	Image string `xorm:"index varchar(1000) not null"`
}

type UserMbox struct {
	User   int64  `xorm:"index"`
	Local  string `xorm:"index varchar(1000) not null"`
	Domain string `xorm:"index varchar(1000) not null"`
}

var (
	db *xorm.Engine
)

func init() {
	var err error

	db, err = xorm.NewEngine("mysql", *dsn+"?charset=utf8&parseTime=True&autocommit=true")
	if err != nil {
		log.Panicln(err)
	}
	db.SetMaxIdleConns(4)
	if *debug {
		db.ShowSQL = true
	}
}
