package controller

import (
	"dousheng/common"
	"dousheng/model"
	"dousheng/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type FavoriteListRequest struct {
	UserID int64  `json:"user_id" form:"user_id" binding:"required"`
	Token  string `json:"token" form:"token" binding:"required"`
}

type FavoriteListResponse struct {
	common.Response
	VideoList []model.Video `json:"video_list,omitempty" form:"video_list"`
}

// 解决接口不一致问题

type FavoriteActionRequest struct {
	//UserID     int64  `json:"user_id" form:"user_id" binding:"required"`
	Token      string `json:"token" form:"token" binding:"required"`
	VideoID    int64  `form:"video_id" json:"video_id" binding:"required"`
	ActionType int32  `json:"action_type" form:"action_type" binding:"required"`
}

type FavoriteActionResponse struct {
	common.Response
}

func FavoriteAction(ctx *gin.Context) {
	var request FavoriteActionRequest
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
	// 登录用户名不匹配
	if u.Username != username {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserNotExisted,
			StatusMsg:  "Login Token error",
		})
		return
	}

	var err error
	if request.ActionType == 1 {
		err = service.Favorite(u.ID, request.VideoID)
	} else if request.ActionType == 2 {
		err = service.UnFavorite(u.ID, request.VideoID)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.OperationFailed,
			StatusMsg:  "Favorite Operation Failed",
		})
		return
	}
	ctx.JSON(http.StatusOK, common.Response{
		StatusCode: common.OK,
		StatusMsg:  "Success",
	})
}

func FavoriteList(ctx *gin.Context) {
	var request FavoriteListRequest
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

	// 获取所有点赞视频
	if err, videoList := service.GetFavoriteListByUserID(u.ID); err != nil {
		ctx.JSON(http.StatusOK, FavoriteListResponse{
			Response: common.Response{
				StatusCode: common.OperationFailed,
				StatusMsg:  "Favorite List OperationFailed",
			},
			VideoList: nil,
		})
	} else {
		ctx.JSON(http.StatusOK, FavoriteListResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "Success",
			},
			VideoList: videoList,
		})
	}
}
