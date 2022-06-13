package service

import (
	"dousheng/database"
	"dousheng/model"
	"log"
)

func CreateVideo(video *model.Video) {
	res := database.MySQLDB.Create(&video)
	if res.Error != nil {
		log.Println("Insert video failed!" + res.Error.Error())
	}

}
func GetVideoById(Id int64) (video *model.Video) {
	var v model.Video
	res := database.MySQLDB.Model(&model.Video{}).Where("id = ?", Id).First(&v)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}
	return &v
}

func GetVideoListByTime(latestTime int64, limit int) (error, []model.Video) {
	var videoList []model.Video
	res := database.MySQLDB.Model(&model.Video{}).Where("publish_time < ?", latestTime).Order("publish_time DESC").Limit(limit).Find(&videoList)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}
	return res.Error, videoList
}

func GetVideoListByUserID(userID int64) (error, []model.Video) {
	var videoList []model.Video
	res := database.MySQLDB.Model(&model.Video{}).Where("author_id = ?", userID).Find(&videoList)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}
	return res.Error, videoList
}
