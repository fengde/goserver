package serviceUser

import (
	"context"
	"fmt"
	"goserver/global"
	"goserver/model"

	"github.com/fengde/gocommon/cachex/localcachex"
	"github.com/fengde/gocommon/errorx"
	"github.com/fengde/gocommon/jsonx"
	"gorm.io/gorm"
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
	Super        int64
}

func GetUserInfo(ctx context.Context, userId int64) (*UserInfo, error) {
	var user model.User
	if err := global.DB.Preload("UserRole").Where("id=?", userId).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, errorx.WithStack(err)
	}

	var roleIds []int64
	for _, ur := range user.UserRole {
		roleIds = append(roleIds, ur.Id)
	}

	return &UserInfo{
		Id:           userId,
		Name:         user.Name,
		RegisterTime: user.CreateTime.String(),
		RoleIds:      roleIds,
		Super:        user.Super,
	}, nil
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

// key=user_id, value=isSuper
var superCache = localcachex.NewLocalCache()

// 是否为超管
func IsSuper(ctx context.Context, userId int64) bool {
	if userId < 1 {
		return false
	}

	isSuper, err := superCache.Get(jsonx.MarshalNoErr(userId))
	if err != nil {
		user, err := GetUserInfo(ctx, userId)
		if err != nil {
			return false
		}
		superCache.Set(jsonx.MarshalNoErr(userId), jsonx.MarshalNoErr(user.Super))

		isSuper = jsonx.MarshalNoErr(user.Super)
	}

	return string(isSuper) == "1"
}
