package serviceUser

import (
	"context"
	"fmt"
	"goserver/global"
	"goserver/model"

	"github.com/fengde/gocommon/errorx"
	"github.com/fengde/gocommon/logx"
)

type Role struct {
	Id   int64
	Name string
}

// 获取用户角色
func GetUserRoles(userId int64) ([]Role, error) {

	var rows []model.UserRole
	if err := global.DB.Preload("Role").Where("user_id=?", userId).Find(&rows).Error; err != nil {
		return nil, errorx.WithStack(err)
	}

	var roles []Role
	for _, row := range rows {
		roles = append(roles, Role{
			Id:   row.Role.Id,
			Name: row.Role.Name,
		})
	}

	return roles, nil
}

type UserInfo struct {
	Id           int64
	Name         string
	RoleIds      []int64
	RegisterTime string
}

func GetUserInfo(userId int64) UserInfo {
	var user model.User
	if err := global.DB.Preload("UserRole").Where("id=?", userId).First(&user).Error; err != nil {
		logx.Error(err)
		return UserInfo{}
	}

	var roleIds []int64
	for _, ur := range user.UserRole {
		roleIds = append(roleIds, ur.Id)
	}

	return UserInfo{
		Id:           userId,
		Name:         user.Name,
		RegisterTime: user.CreateTime.String(),
		RoleIds:      roleIds,
	}
}

// 绑定角色
func BindRole(ctx context.Context, userIds []int64, roleIds []int64, createUserId int64) error {
	if len(userIds) == 0 || len(roleIds) == 0 {
		return nil
	}

	userRoles := []model.UserRole{}
	for _, userId := range userIds {
		for _, roleId := range roleIds {
			userRoles = append(userRoles, model.UserRole{
				UserId:       userId,
				RoleId:       roleId,
				CreateUserId: createUserId,
			})
		}
	}

	if err := global.DB.Create(&userRoles).Error; err != nil {
		return errorx.WithStack(err)
	}

	var roles []string
	for _, roleId := range roleIds {
		roles = append(roles, fmt.Sprintf("role_%d", roleId))
	}
	// 用户与角色关系 加入casbin
	for _, userId := range userIds {
		if _, err := global.Enforcer.AddRolesForUser(fmt.Sprintf("user_%d", userId), roles); err != nil {
			return errorx.WithStack(err)
		}
	}
	return nil
}
