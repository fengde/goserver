package handler

import (
	"goserver/service/serviceRbac"

	"github.com/fengde/gocommon/errorx"
)

type GetPermissionGroupsRequest struct {
}

type GetPermissionGroupsResponse struct {
	PermissionGroups []serviceRbac.PermissionGroup `json:"permission_groups"`
}

// 查询所有权限组
func GetPermissionGroups(c *Context, request *GetPermissionGroupsRequest) (*GetPermissionGroupsResponse, error) {
	groups, err := serviceRbac.GetPermissionGroups(c.GetCtx())
	if err != nil {
		return nil, errorx.WithStack(err)
	}

	return &GetPermissionGroupsResponse{
		PermissionGroups: groups,
	}, nil
}
