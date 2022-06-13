package controller

import (
	"dousheng/common"
	"dousheng/model"
	"dousheng/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RelationActionRequest struct {
	Token      string `json:"token" form:"token" binding:"required"`
	ToUserID   int64  `json:"to_user_id" form:"to_user_id" binding:"required"`
	ActionType int32  `json:"action_type" form:"action_type" binding:"required,oneof=1 2"`
}

type RelationRequest struct {
	UserID int64  `json:"user_id" form:"user_id" binding:"required"`
	Token  string `json:"token" form:"token" binding:"required"`
}

type RelationListResponse struct {
	common.Response
	UserList []model.UserInfo `json:"user_list,omitempty" binding:"required"`
}

func RelationAction(ctx *gin.Context) {
	var request RelationActionRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}
	// 判断用户登录
	strs := strings.Split(request.Token, ":")
	username := strs[0]
	u, _ := service.GetUserByUsername(username)
	if u.Username != username {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserNotExisted,
			StatusMsg:  "Login Token error",
		})
		return
	}
	// 关注取关操作 没做验证 比如判断被关注的是否存在 能否自关
	var err error
	if request.ActionType == 1 {
		err = service.Follow(u.ID, request.ToUserID)
	} else if request.ActionType == 2 {
		err = service.UnFollow(u.ID, request.ToUserID)
	}
	if err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.OperationFailed,
			StatusMsg:  "Operation Failed",
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Response{
		StatusCode: common.OK,
		StatusMsg:  "Success",
	})
}

func FollowList(ctx *gin.Context) {
	var request RelationRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}
	// 判断用户登录
	strs := strings.Split(request.Token, ":")
	username := strs[0]
	u, _ := service.GetUserByUsername(username)
	if u.Username != username {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserNotExisted,
			StatusMsg:  "Login Token error",
		})
		return
	}
	if err, userList := service.GetFollowListByUserID(u.ID, request.UserID); err != nil {
		ctx.JSON(http.StatusOK, RelationListResponse{
			Response: common.Response{
				StatusCode: common.OperationFailed,
				StatusMsg:  "OperationFailed",
			},
			UserList: nil,
		})
	} else {
		ctx.JSON(http.StatusOK, RelationListResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "Success",
			},
			UserList: userList,
		})
	}
}

func FollowerList(ctx *gin.Context) {
	var request RelationRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}
	// 判断用户登录
	strs := strings.Split(request.Token, ":")
	username := strs[0]
	u, _ := service.GetUserByUsername(username)
	if u.Username != username {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserNotExisted,
			StatusMsg:  "Login Token error",
		})
		return
	}

	if err, userList := service.GetFollowerListByUserID(u.ID, request.UserID); err != nil {
		ctx.JSON(http.StatusOK, RelationListResponse{
			Response: common.Response{
				StatusCode: common.OperationFailed,
				StatusMsg:  "OperationFailed",
			},
			UserList: nil,
		})
	} else {
		ctx.JSON(http.StatusOK, RelationListResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "Success",
			},
			UserList: userList,
		})
	}
}
