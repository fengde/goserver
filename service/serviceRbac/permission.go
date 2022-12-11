package serviceRbac

import (
	"context"
	"goserver/global"
	"goserver/model"

	"github.com/fengde/gocommon/errorx"
)

type PermissionUrl struct {
	PermissionId int64
	Url          string
	Method       string
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
