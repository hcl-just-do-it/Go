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

type FeedRequest struct {
	LatestTime int64  `json:"latest_time,omitempty"`
	Token      string `json:"token,omitempty"`
}

type VideoResponse struct {
	Id            int64          `json:"id,omitempty"`
	Author        model.UserInfo `json:"author"`
	PlayUrl       string         `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string         `json:"cover_url,omitempty"`
	FavoriteCount int64          `json:"favorite_count,omitempty"`
	CommentCount  int64          `json:"comment_count,omitempty"`
	Title         string         `json:"title,omitempty"`
}

type FeedResponse struct {
	common.Response
	VideoList []VideoResponse `json:"video_list"`
	NextTime  int64           `json:"next_time,omitempty"`
}

const LIMIT = 30

// Feed same demo video list for every request
func Feed(ctx *gin.Context) {
	var request FeedRequest
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "Parameter parsing error",
		})
		return
	}
	// 如果没传时间值 那么默认是当前时间
	if request.LatestTime == 0 {
		request.LatestTime = time.Now().Unix()
	}
	// 查videoList
	if err, videoList := service.GetVideoListByTime(request.LatestTime, LIMIT); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.OperationFailed,
			StatusMsg:  "OperationFailed",
		})
	} else {
		// 查最后一个video
		index := len(videoList) - 1
		var nextTime int64
		if index >= 0 {
			nextTime = videoList[index].PublishTime
		} else {
			nextTime = time.Now().Unix()
		}
		//var ids = make([]int64, len(videoList))
		//for i, v := range videoList {
		//	ids[i] = v.AuthorID
		//}
		// 判断是否有token
		var loginID int64
		loginID = 0
		if request.Token != "" {
			username := strings.Split(request.Token, ":")[0]
			u, _ := service.GetUserByUsername(username)
			loginID = u.ID
		}
		//_, userInfoList := service.GetUserInfoListByIDs(loginID, ids)
		var responseList = make([]VideoResponse, len(videoList))
		for i, v := range videoList {
			_, userInfo := service.GetUserInfoByUserID(loginID, v.AuthorID)
			responseList[i] = VideoResponse{
				Id: v.Id,
				//Author:        userInfoList[i],
				Author:        userInfo,
				PlayUrl:       v.PlayUrl,
				CoverUrl:      v.CoverUrl,
				FavoriteCount: v.FavoriteCount,
				CommentCount:  v.CommentCount,
				Title:         v.Title,
			}
		}
		ctx.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "Success",
			},
			VideoList: responseList,
			NextTime:  nextTime,
		})
	}

}
