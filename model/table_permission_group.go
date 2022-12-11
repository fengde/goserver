package model

import (
	"time"
)

type PermissionGroup struct {
	Id          int64        `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	Name        string       `gorm:"column:name;default:;NOT NULL;comment:'权限组名称'"`
	CreateTime  time.Time    `gorm:"column:create_time;default:NULL"`
	UpdateTime  time.Time    `gorm:"column:update_time;default:NULL"`
	Permissions []Permission `gorm:"foreignKey:permission_group_id"`
}

func (p *PermissionGroup) TableName() string {
	return "permission_group"
}
