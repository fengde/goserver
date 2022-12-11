package serviceRbac

import (
	"context"
	"fmt"
	"goserver/global"
	"goserver/model"
	"time"

	"github.com/fengde/gocommon/errorx"
	"github.com/fengde/gocommon/logx"
	"github.com/fengde/gocommon/taskx"
	"gorm.io/gorm"
)

type NewRoleParams struct {
	Name          string
	PermissionIds []int64
	CreateUserId  int64
}

// 创建角色
func NewRole(ctx context.Context, params NewRoleParams) (int64, error) {
	exist, err := GetRoleByName(ctx, params.Name)

	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, errorx.WithStack(err)
	}
	if exist != nil {
		return 0, errorx.New("角色名称已经存在")
	}

	var newRoleId int64
	current := time.Now()
	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		role := model.Role{
			Name:         params.Name,
			CreateUserId: params.CreateUserId,
			CreateTime:   current,
			UpdateTime:   current,
		}
		if err := tx.Create(&role).Error; err != nil {
			return errorx.WithStack(err)
		}

		if role.Id > 0 {
			rolePermissions := []model.RolePermission{}
			for _, permissionId := range params.PermissionIds {
				rolePermissions = append(rolePermissions, model.RolePermission{
					RoleId:       role.Id,
					PermissionId: permissionId,
					CreateTime:   current,
					UpdateTime:   current,
				})
			}
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return errorx.WithStack(err)
			}
		}

		newRoleId = role.Id

		return nil
	}); err != nil {
		return 0, errorx.WithStack(err)
	}

	// 角色与permission url加入casbin
	permissionUrls, err := GetPermissionUrls(ctx, params.PermissionIds)
	if err != nil {
		return 0, errorx.WithStack(err)
	}

	var rules [][]string
	for _, permissionUrl := range permissionUrls {
		rules = append(rules, []string{fmt.Sprintf("role_%d", newRoleId), permissionUrl.Url, permissionUrl.Method})
	}

	if _, err := global.Enforcer.AddPolicies(rules); err != nil {
		return 0, errorx.WithStack(err)
	}

	return newRoleId, nil
}

type Role struct {
	Id   int64
	Name string
}

// 根据名称获取角色信息
func GetRoleByName(ctx context.Context, name string) (*Role, error) {
	var row model.Role
	if err := global.DB.Where("name = ?", name).First(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, errorx.WithStack(err)
	}

	if row.Id == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &Role{
		Id:   row.Id,
		Name: row.Name,
	}, nil
}

// 删除角色
func DeleteRole(ctx context.Context, roleIds []int64) error {
	if len(roleIds) == 0 {
		return nil
	}

	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		return taskx.NewSerialTaskGroup(func() error {
			return tx.Where("role_id IN ?", roleIds).Delete(model.RolePermission{}).Error
		}, func() error {
			return tx.Where("role_id IN ?", roleIds).Delete(model.UserRole{}).Error
		}, func() error {
			return tx.Where("id IN ?", roleIds).Delete(model.Role{}).Error
		}).Run()
	}); err != nil {
		return errorx.WithStack(err)
	}

	// 移除casbin role相关信息
	for _, roleId := range roleIds {
		_, err := global.Enforcer.DeleteRole(fmt.Sprintf("role_%d", roleId))
		if err != nil {
			logx.ErrorWithCtx(ctx, err)
		}
	}

	return nil
}
