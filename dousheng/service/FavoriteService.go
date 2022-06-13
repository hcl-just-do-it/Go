package service

import (
	"dousheng/database"
	"dousheng/model"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

func Favorite(fromID int64, vID int64) error {
	favoriteKey := database.GetFavoriteKey(fromID)
	videoKey := database.GetVideoKey(vID)

	score := float64(time.Now().Unix())

	// redis事务
	pipe := database.RedisDB.TxPipeline()
	defer pipe.Close()

	// 表明用户喜欢哪些视频
	pipe.ZAdd(favoriteKey, redis.Z{Score: score, Member: vID})
	// 表明视频是否被某用户喜欢
	pipe.ZAdd(videoKey, redis.Z{Score: score, Member: fromID})

	_, err := pipe.Exec()
	if err != nil {
		pipe.Discard()
		return err
	}

	return nil
}

func UnFavorite(fromID int64, vID int64) error {
	favoriteKey := database.GetFavoriteKey(fromID)
	videoKey := database.GetVideoKey(vID)

	// redis事务
	pipe := database.RedisDB.TxPipeline()
	defer pipe.Close()

	// 用户不再喜欢该视频，该视频也不再被用户喜欢
	pipe.ZRem(favoriteKey, vID)
	pipe.ZRem(videoKey, fromID)

	_, err := pipe.Exec()
	if err != nil {
		pipe.Discard()
		return err
	}

	return nil
}

// 通过用户ID获取他喜爱的视频列表

func GetFavoriteListByUserID(userID int64) (error, []model.Video) {
	favoriteKey := database.GetFavoriteKey(userID)
	var videoInfoList []model.Video
	var err error
	var idstrs []string

	// 返回该用户喜欢视频的id
	idstrs, err = database.RedisDB.ZRange(favoriteKey, 0, -1).Result()
	ids := make([]int64, len(idstrs))
	for i, v := range idstrs {
		ids[i], _ = strconv.ParseInt(v, 10, 64)
	}
	err, videoInfoList = GetVideoInfoListByIDs(ids)
	return err, videoInfoList
}

// 根据视频id返回视频信息切片

func GetVideoInfoListByIDs(ids []int64) (error, []model.Video) {
	var videos []model.Video
	var err error
	res := database.MySQLDB.Model(&model.Video{}).Where("id in (?)", ids).Find(&videos)
	if res.Error != nil {
		log.Println("query userinfo by ids failed!" + res.Error.Error())
		return res.Error, nil
	}
	return err, videos
}
