package model

import (
	"time"
)

type UserRole struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	UserId       int64     `gorm:"column:user_id;default:0;NOT NULL"`
	RoleId       int64     `gorm:"column:role_id;default:0;NOT NULL"`
	CreateUserId int64     `gorm:"column:create_user_id;default:0;NOT NULL"`
	CreateTime   time.Time `gorm:"column:create_time;default:NULL"`
	UpdateTime   time.Time `gorm:"column:update_time;default:NULL"`
	Role         Role      `gorm:"foreignKey:role_id"`
	User         User      `gorm:"foreignKey:user_id"`
}

func (u *UserRole) TableName() string {
	return "user_role"
}
