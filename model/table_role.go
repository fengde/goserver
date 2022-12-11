package model

import (
	"time"
)

type Role struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name         string    `gorm:"column:name;default:;NOT NULL;comment:'角色名称'"`
	CreateUserId int64     `gorm:"column:create_user_id;default:0;NOT NULL"`
	UpdateUserId int64     `gorm:"column:update_user_id;default:0;NOT NULL"`
	CreateTime   time.Time `gorm:"column:create_time;default:NULL"`
	UpdateTime   time.Time `gorm:"column:update_time;default:NULL"`
}

func (r *Role) TableName() string {
	return "role"
}
