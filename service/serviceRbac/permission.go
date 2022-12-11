package serviceRbac

import (
	"context"
	"goserver/global"
	"goserver/model"

	"github.com/fengde/gocommon/errorx"
	"gorm.io/gorm"
)

type PermissionUrl struct {
	PermissionId int64  `json:"id"`
	Url          string `json:"url"`
	Method       string `json:"method"`
}

// 查下权限对应的url, method
func GetPermissionUrls(ctx context.Context, permissionIds []int64) ([]PermissionUrl, error) {
	var rows []model.PermissionUrl
	if err := global.DB.Where("permission_id in ?", permissionIds).Find(&rows).Error; err != nil {
		return nil, errorx.WithStack(err)
	}
	var back []PermissionUrl
	for _, row := range rows {
		back = append(back, PermissionUrl{
			PermissionId: row.PermissionId,
			Url:          row.Url,
			Method:       row.Method,
		})
	}
	return back, nil
}

type Permission struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type PermissionGroup struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
}

// 返回权限组详情
func GetPermissionGroups(ctx context.Context) ([]PermissionGroup, error) {

	var rows []model.PermissionGroup
	if err := global.DB.Preload("Permissions").Find(&rows).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, errorx.WithStack(err)
		}
		return nil, nil
	}

	var pgs []PermissionGroup

	for _, row := range rows {
		var permissions []Permission
		for i := range row.Permissions {
			permissions = append(permissions, Permission{
				Id:   row.Permissions[i].Id,
				Name: row.Permissions[i].Name,
			})
		}
		pgs = append(pgs, PermissionGroup{
			Id:          row.Id,
			Name:        row.Name,
			Permissions: permissions,
		})
	}

	return pgs, nil
}
