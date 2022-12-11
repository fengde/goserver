package model

import (
	"time"
)

type User struct {
	Id         int64      `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name       string     `gorm:"column:name;default:;NOT NULL"`
	Email      string     `gorm:"column:email;default:;NOT NULL"`
	Password   string     `gorm:"column:password;NOT NULL"`
	Super      int64      `gorm:"column:super;default:0;NOT NULL;comment:'是否超管'"`
	CreateTime time.Time  `gorm:"column:create_time;default:NULL"`
	UpdateTime time.Time  `gorm:"column:update_time;default:NULL"`
	UserRole   []UserRole `gorm:"foreignKey:user_id"`
}

func (u *User) TableName() string {
	return "user"
}
