package controller

import (
	"dousheng/common"
	"dousheng/model"
	"dousheng/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type CommentActionRequest struct {
	Token       string `form:"token" json:"token" binding:"required"`
	VideoID     int64  `form:"video_id" json:"video_id" binding:"required"`
	ActionType  int32  `form:"action_type" json:"action_type" binding:"required"`
	CommentText string `form:"comment_text" json:"comment_text" `
	CommentID   int64  `form:"comment_id" json:"comment_id" `
}

type CommentActionResponse struct {
	common.Response
	Comment model.CommentInfo `json:"comment"`
}

func CommentAction(c *gin.Context) {
	var request CommentActionRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}

	choose := request.ActionType
	if choose == 1 { // 发布评论
		// 获取userid
		strs := strings.Split(request.Token, ":")
		username := strs[0]
		u, _ := service.GetUserByUsername(username)
		if u.Username != username {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: common.UserNotExisted,
				StatusMsg:  "Login Token error",
			})
			return
		}

		// 获取userinfo
		_, newUserInfo := service.GetUserInfoByUserID(u.ID, u.ID)

		newComment := model.Comment{
			VideoID:    request.VideoID,
			Content:    request.CommentText,
			CreateDate: time.Now().Format("01-02"),
			UserID:     u.ID,
		}

		service.CreateComment(&newComment)

		newCommentInfo := model.CommentInfo{
			VideoID:    request.VideoID,
			Content:    request.CommentText,
			CreateDate: newComment.CreateDate,
			User:       newUserInfo,
			CommentID:  newComment.CommentID,
		}

		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "success",
			},
			Comment: newCommentInfo,
		})
	} else if choose == 2 { // 删除评论
		service.DeleteComment(request.CommentID)
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "success",
			},
		})
	}

}

type CommentListRequest struct {
	Token   string `form:"token" json:"token" binding:"required"`
	VideoID int64  `form:"video_id" json:"video_id" binding:"required"`
}

type CommentListResponse struct {
	common.Response
	CommentList []model.CommentInfo `json:"comment_list"`
}

func CommentList(c *gin.Context) {
	var request CommentListRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}
	// 获取userid
	strs := strings.Split(request.Token, ":")
	username := strs[0]
	u, _ := service.GetUserByUsername(username)
	if u.Username != username {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserNotExisted,
			StatusMsg:  "Login Token error",
		})
		return
	}
	vid := request.VideoID
	newCommentList := service.GetCommentByVideoID(u.ID, vid)

	c.JSON(http.StatusOK, CommentListResponse{
		Response: common.Response{
			StatusCode: common.OK,
			StatusMsg:  "success",
		},
		CommentList: newCommentList,
	})
}
