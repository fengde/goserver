package model

import (
	"time"
)

type PermissionUrl struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	PermissionId int64     `gorm:"column:permission_id;default:0;NOT NULL;comment:'权限id'"`
	Url          string    `gorm:"column:url;default:;NOT NULL;comment:'http url'"`
	Method       string    `gorm:"column:method;default:;NOT NULL;comment:'http method'"`
	CreateTime   time.Time `gorm:"column:create_time;default:NULL"`
	UpdateTime   time.Time `gorm:"column:update_time;default:NULL"`
}

func (p *PermissionUrl) TableName() string {
	return "permission_url"
}
