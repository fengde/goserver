package model

type Permission struct {
	Id                int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL"`
	PermissionGroupId string `gorm:"column:permission_group_id;default:0;NOT NULL;comment:'权限组id名称'"`
	Name              string `gorm:"column:name;default:;NOT NULL;comment:'权限名称'"`
}

func (p *Permission) TableName() string {
	return "permission"
}
