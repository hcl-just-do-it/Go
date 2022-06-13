package service

import (
	"dousheng/database"
	"dousheng/model"
	"github.com/go-redis/redis"
	"log"
	"strconv"
	"time"
)

// 通过用户名返回UserInfo
func GetUserByUsername(username string) (user *model.User, err error) {
	var u model.User
	res := database.MySQLDB.Model(&model.User{}).Where("username = ?", username).First(&u)
	if res.Error != nil {
		log.Println(res.Error.Error())
	}
	return &u, res.Error
}

func IsUserExisted(username string) bool {
	var user model.User
	res := database.MySQLDB.Where("username=?", username).First(&user)
	if res.RowsAffected >= 1 {
		return true
	}
	return false
}

func CreateUser(user *model.User) {
	res := database.MySQLDB.Create(&user)
	if res.Error != nil {
		log.Println("Insert user failed!" + res.Error.Error())
	}
	//return user
}

func Follow(fromID int64, toID int64) error {
	followKey := database.GetFollowKey(fromID)
	followerKey := database.GetFollowerKey(toID)
	score := float64(time.Now().Unix())

	// redis事务
	pipe := database.RedisDB.TxPipeline()
	defer pipe.Close()

	pipe.ZAdd(followKey, redis.Z{Score: score, Member: toID})
	pipe.ZAdd(followerKey, redis.Z{Score: score, Member: fromID})

	_, err := pipe.Exec()
	if err != nil {
		pipe.Discard()
		return err
	}

	return nil
}

func UnFollow(fromID int64, toID int64) error {
	followKey := database.GetFollowKey(fromID)
	followerKey := database.GetFollowerKey(toID)

	// redis事务
	pipe := database.RedisDB.TxPipeline()
	defer pipe.Close()

	pipe.ZRem(followKey, toID)
	pipe.ZRem(followerKey, fromID)

	_, err := pipe.Exec()
	if err != nil {
		pipe.Discard()
		return err
	}

	return nil
}

func GetFollowListByUserID(loginID int64, userID int64) (error, []model.UserInfo) {
	followKey := database.GetFollowKey(userID)
	var userInfoList []model.UserInfo
	var err error
	var idstrs []string

	idstrs, err = database.RedisDB.ZRange(followKey, 0, -1).Result()
	ids := make([]int64, len(idstrs))
	for i, v := range idstrs {
		ids[i], _ = strconv.ParseInt(v, 10, 64)
	}
	err, userInfoList = GetUserInfoListByIDs(loginID, ids)
	return err, userInfoList
}

func GetFollowerListByUserID(loginID int64, userID int64) (error, []model.UserInfo) {
	followerKey := database.GetFollowerKey(userID)
	var userInfoList []model.UserInfo
	var err error
	var idstrs []string

	idstrs, err = database.RedisDB.ZRange(followerKey, 0, -1).Result()
	ids := make([]int64, len(idstrs))
	for i, v := range idstrs {
		ids[i], _ = strconv.ParseInt(v, 10, 64)
	}
	err, userInfoList = GetUserInfoListByIDs(loginID, ids)
	return err, userInfoList
}

func GetFollowCount(userID int64) (int64, error) {
	followKey := database.GetFollowKey(userID)
	return database.RedisDB.ZCard(followKey).Result()
}

func GetFollowerCount(userID int64) (int64, error) {
	followerKey := database.GetFollowerKey(userID)
	return database.RedisDB.ZCard(followerKey).Result()
}

func IsFollow(fromID int64, toID int64) bool {
	followKey := database.GetFollowKey(fromID)
	score, _ := database.RedisDB.ZScore(followKey, strconv.FormatInt(toID, 10)).Result()
	if score > 0 {
		return true
	} else {
		return false
	}
}

func GetUserInfoListByIDs(loginID int64, ids []int64) (error, []model.UserInfo) {
	var users []model.User
	var err error
	res := database.MySQLDB.Model(&model.User{}).Where("id in (?)", ids).Find(&users)
	if res.Error != nil {
		log.Println("query userinfo by ids failed!" + res.Error.Error())
		return res.Error, nil
	}
	userInfos := make([]model.UserInfo, len(users))
	for i, v := range users {
		var followCount, followerCount int64
		var isFollow bool
		followCount, err = GetFollowCount(v.ID)
		followerCount, err = GetFollowerCount(v.ID)
		isFollow = IsFollow(loginID, v.ID)
		userInfos[i] = model.UserInfo{
			Id:            v.ID,
			Name:          v.Username,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      isFollow,
		}
	}
	return err, userInfos
}

func GetUserInfoByUserID(loginID int64, userID int64) (error, model.UserInfo) {
	var user model.User
	var userInfo model.UserInfo
	var err error
	res := database.MySQLDB.Model(&model.User{}).Where("id = ?", userID).First(&user)
	if res.Error != nil {
		log.Println("query userinfo by id failed!" + res.Error.Error())
		return res.Error, userInfo
	}
	var followCount, followerCount int64
	var isFollow bool
	followCount, err = GetFollowCount(userID)
	followerCount, err = GetFollowerCount(userID)
	isFollow = IsFollow(loginID, userID)
	userInfo = model.UserInfo{
		Id:            userID,
		Name:          user.Username,
		FollowCount:   followCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}
	return err, userInfo
}
