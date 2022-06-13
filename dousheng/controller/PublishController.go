package controller

import (
	"dousheng/common"
	"dousheng/model"
	"dousheng/service"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	idworker "github.com/gitstliu/go-id-worker"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"golang.org/x/net/context"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PublishRequest struct {
	Token string `form:"token" json:"token" binding:"required"`
	Data  []byte `form:"data" json:"data" binding:"required"`
	Title string `form:"title" json:"title" binding:"required"`
}

type PublishListResponse struct {
	common.Response
	VideoList []VideoResponse `json:"video_list"`
}

type PublishResponse struct {
	common.Response
}

var domainName = "http://rd5met9ed.hn-bkt.clouddn.com"
var bucket = "top-20"

//var domainName = "qiniu.jianggang.top"
//var bucket = "haiwai-20"
var accessKey = "ANvRMQN-FX6C6abeKAYxqAq1qq9je2x1UAmlLjFA"
var secretKey = "RhH86hgmwDphJxs5jBa1yUzZM7ydAch7msd-_VSi"
var videoFileExt = []string{"mp4", "flv"} //此处可根据需要添加格式
var idGen *idworker.IdWorker

func init() {
	idGen = &idworker.IdWorker{}
	idGen.InitIdWorker(1, 1)
}

func Publish(c *gin.Context) {
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "FormFile parsing  error",
		})
		return
	}
	errd, coverUrl, playUrl := UploadVideo(file)
	if errd != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: common.ParamInvalid,
			StatusMsg:  "UploadVideo  error",
		})
		return
	}
	token, _ := c.GetPostForm("token")
	strs := strings.Split(token, ":")
	username := strs[0]
	title, _ := c.GetPostForm("title")
	u, _ := service.GetUserByUsername(username)

	//存数据库之后，再根据数据库隐式生成的ID，存入video的ID
	video := model.Video{
		AuthorID:      u.ID,
		PlayUrl:       domainName + "/" + playUrl,
		CoverUrl:      domainName + "/" + coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		PublishTime:   time.Now().Unix(),
	}
	fmt.Println("video=%#v\n", video)
	service.CreateVideo(&video)
	c.JSON(http.StatusOK, common.Response{
		StatusCode: common.OK,
		StatusMsg:  "Publish Success!",
	})
}

func IsVideoAllowed(suffix string) bool {
	for _, fileExt := range videoFileExt {
		if suffix == fileExt {
			return true
		}
	}
	return false
}

func UploadVideo(file *multipart.FileHeader) (err error, coverUrl string, playUrl string) {
	//先处理输入
	filename := file.Filename                      //获取文件名
	indexOfDot := strings.LastIndex(filename, ".") //获取文件后缀名前的.的位置
	if indexOfDot < 0 {
		return errors.New("没有获取到文件的后缀名"), coverUrl, playUrl
	}
	suffix := filename[indexOfDot+1 : len(filename)] //获取后缀名
	suffix = strings.ToLower(suffix)                 //后缀名统一小写处理
	if !IsVideoAllowed(suffix) {
		return errors.New("上传的文件不符合视频的格式"), coverUrl, playUrl
	}
	fmt.Println("刚才上传的文件后缀名：" + suffix)
	id, err := idGen.NextId()
	filename = strconv.FormatInt(id, 10)

	videoName := filename + "." + suffix //视频的文件名
	data, err := file.Open()             //data是文件内容的访问接口（重点）
	VideoFolderName := "video"
	key := VideoFolderName + "/" + videoName //key是要上传的文件访问路径
	//下面是七牛api
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	//视频封面start
	coverName := filename + "." + "jpg"           //封面的文件名
	coverFolderName := "cover"                    //七牛云中存放图片的目录名。用于与文件名拼接，组成文件路径
	coverKey := coverFolderName + "/" + coverName //封面的访问路径，我们通过此路径在七牛云空间中定位封面
	saveJpgEntry := base64.StdEncoding.EncodeToString([]byte(bucket + ":" + coverKey))
	putPolicy.PersistentOps = "vframe/jpg/offset/1/w/534/h/949|saveas/" + saveJpgEntry //取视频第1秒的截图
	//   "vframe/jpg/offset/7/w/480/h/360",
	//end
	putPolicy.Expires = 7200 //自定义凭证有效期（示例2小时，Expires 单位为秒，为上传凭证的有效时间）
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	//data是字节流，data := []byte("hello, this is qiniu cloud")
	//file.size是要上传的文件大小
	err = formUploader.Put(context.Background(), &ret, upToken, key, data, file.Size, &putExtra)
	if err != nil {
		return err, coverUrl, playUrl
	}
	//fmt.Println(ret.Key, ret.Hash)

	coverUrl = coverKey
	playUrl = key
	return err, coverUrl, playUrl
}

func PublishList(ctx *gin.Context) { //我发布的视频列表
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
	if err, videoList := service.GetVideoListByUserID(u.ID); err != nil {
		ctx.JSON(http.StatusOK, common.Response{
			StatusCode: common.OperationFailed,
			StatusMsg:  "OperationFailed",
		})
	} else {
		var userInfo model.UserInfo
		var responseList = make([]VideoResponse, len(videoList))
		for i, v := range videoList {
			if i == 0 {
				_, userInfo = service.GetUserInfoByUserID(u.ID, v.AuthorID)
			}
			responseList[i] = VideoResponse{
				Id:            v.Id,
				Author:        userInfo,
				PlayUrl:       v.PlayUrl,
				CoverUrl:      v.CoverUrl,
				FavoriteCount: v.FavoriteCount,
				CommentCount:  v.CommentCount,
				Title:         v.Title,
			}
		}
		ctx.JSON(http.StatusOK, PublishListResponse{
			Response: common.Response{
				StatusCode: common.OK,
				StatusMsg:  "Success",
			},
			VideoList: responseList,
		})
	}
}
