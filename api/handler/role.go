package handler

import (
	"goserver/service/serviceRbac"
	"goserver/service/serviceUser"

	"github.com/fengde/gocommon/errorx"
)

type NewRoleRequest struct {
	Name          string  `json:"name" valid:"required~不允许为空"`
	PermissionIds []int64 `json:"permission_ids" valid:"required~不允许为空"`
	UserIds       []int64 `json:"user_ids"`
}

type NewRoleResponse struct {
	Id int64 `json:"id"`
}

// 新建角色，同时绑定用户
func NewRole(c *Context, request *NewRoleRequest) (*NewRoleResponse, error) {

	ctx := c.GetCtx()

	roleId, err := serviceRbac.NewRole(ctx, serviceRbac.NewRoleParams{
		Name:          request.Name,
		PermissionIds: request.PermissionIds,
		CreateUserId:  c.UserId,
	})
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	if len(request.UserIds) > 0 && roleId > 0 {
		if err := serviceUser.BindRole(ctx, request.UserIds, []int64{roleId}, c.UserId); err != nil {
			return nil, errorx.WithStack(err)
		}
	}

	return &NewRoleResponse{
		Id: roleId,
	}, nil
}

type DeleteRoleRequest struct {
	RoleIds []int64 `json:"role_ids" valid:"required~不允许为空"`
	Force   int64
}

// 删除角色
func DeleteRole(c *Context, request *DeleteRoleRequest) error {
	switch request.Force {
	case 0:
		// 检查角色是否用户关联
		// todo
		return nil
	default:
		return serviceRbac.DeleteRole(c.GetCtx(), request.RoleIds)
	}
}
