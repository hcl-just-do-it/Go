package service

import (
	"dousheng/database"
	"dousheng/model"
	"log"
)

func GetCommentByVideoID(login int64, vid int64) []model.CommentInfo {
	var u []model.Comment
	res := database.MySQLDB.Model(&model.Comment{}).Where("video_id = ?", vid).Find(&u)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}

	var ids = make([]int64, len(u))
	for i, v := range u {
		ids[i] = v.UserID
	}

	_, userInfoList := GetUserInfoListByIDs(login, ids)
	var commentInfoList = make([]model.CommentInfo, len(u))
	for i, v := range userInfoList {
		commentInfoList[i] = model.CommentInfo{
			VideoID:    vid,
			Content:    u[i].Content,
			CreateDate: u[i].CreateDate,
			CommentID:  u[i].CommentID,
			User:       v,
		}
	}

	return commentInfoList
}

func CreateComment(comment *model.Comment) {
	res := database.MySQLDB.Create(&comment)
	if res.Error != nil {
		log.Println("Insert user failed!" + res.Error.Error())
	}
	//return user
}

func DeleteComment(cid int64) {
	database.MySQLDB.Delete(&model.Comment{}, cid)
	//return user
}
