package db

import (
	"github.com/fengde/gocommon/storex/mysqlx"
)

type User struct {
	Id         int64             `xorm:"pk autoincr not null int(11) 'id'"`
	Staffid    string            `xorm:"not null varchar(64) 'staffid'"`
	Staffname  string            `xorm:"not null default '' varchar(64) 'staffname'"`
	CreateTime mysqlx.NormalTime `xorm:"not null default 'CURRENT_TIMESTAMP' timestamp 'create_time'"`
	UpdateTime mysqlx.NormalTime `xorm:"not null default 'CURRENT_TIMESTAMP' timestamp 'update_time'"`
}

func (u User) TableName() string {
	return "user"
}
