package controller

import (
	"dousheng/common"
	"dousheng/model"
	"dousheng/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type UserLoginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	common.Response
	User model.UserInfo `json:"user"`
}

func Register(ctx *gin.Context) {
	var request UserLoginRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}
	token := request.Username + ":" + request.Password
	// 参数验证

	// 检查用户是否存在
	if exist := service.IsUserExisted(request.Username); exist {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserHasExisted,
			StatusMsg:  "User has existed",
		})
		return
	}
	//
	newUser := model.User{
		Username: request.Username,
		Password: request.Password,
	}
	service.CreateUser(&newUser)
	ctx.JSON(http.StatusOK, UserLoginResponse{
		Response: common.Response{
			StatusCode: common.OK,
			StatusMsg:  "success",
		},
		UserId: newUser.ID,
		Token:  token,
	})
}

func Login(ctx *gin.Context) {
	var request UserLoginRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}
	// 检查用户是否存在
	u, _ := service.GetUserByUsername(request.Username)
	if u.Username != request.Username {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserNotExisted,
			StatusMsg:  "User not existed",
		})
		return
	}
	// 登录状态保持 session? jwt?
	if u.Password == request.Password {
		ctx.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "success",
			},
			UserId: u.ID,
			Token:  u.Username + ":" + u.Password,
		})
	} else {
		ctx.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{
				StatusCode: common.WrongPassword,
				StatusMsg:  "Wrong password",
			},
		})
	}
}

func User(ctx *gin.Context) {
	token := ctx.Query("token")
	strs := strings.Split(token, ":")
	username := strs[0]
	u, _ := service.GetUserByUsername(username)
	if u.Username != username {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.UserNotExisted,
			StatusMsg:  "User not existed",
		})
		return
	}
	if err, userInfo := service.GetUserInfoByUserID(u.ID, u.ID); err != nil {
		ctx.JSON(http.StatusOK, UserResponse{
			Response: common.Response{
				StatusCode: common.OperationFailed,
			},
			User: userInfo,
		})
	} else {
		ctx.JSON(http.StatusOK, UserResponse{
			Response: common.Response{
				StatusCode: common.OK,
			},
			User: userInfo,
		})
	}

}
