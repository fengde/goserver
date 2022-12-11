package model

import (
	"time"
)

type RolePermission struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	RoleId       int64     `gorm:"column:role_id;default:0;NOT NULL"`
	PermissionId int64     `gorm:"column:permission_id;default:0;NOT NULL"`
	CreateTime   time.Time `gorm:"column:create_time;default:NULL"`
	UpdateTime   time.Time `gorm:"column:update_time;default:NULL"`
}

func (r *RolePermission) TableName() string {
	return "role_permission"
}
