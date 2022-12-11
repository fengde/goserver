package model

import (
	"time"
)

type PermissionWeb struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	PermissionId int64     `gorm:"column:permission_id;default:0;NOT NULL"`
	WebCode      string    `gorm:"column:web_code;default:;NOT NULL;comment:'前端组件编码'"`
	CreateTime   time.Time `gorm:"column:create_time;default:NULL"`
	UpdateTime   time.Time `gorm:"column:update_time;default:NULL"`
}

func (p *PermissionWeb) TableName() string {
	return "permission_web"
}
