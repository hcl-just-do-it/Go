package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/test/", controller.Test)
	apiRouter.POST("/testPost/", controller.TestPOST)
	apiRouter.POST("/testDB2/", controller.TestDB2)
	apiRouter.POST("/testGorm/", controller.TestGorm)
	apiRouter.GET("/feed/", controller.Feed)
	apiRouter.GET("/user/", controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register) //这里是POST请求！
	apiRouter.POST("/user/login/", controller.Login)       //这里是POST请求！
	apiRouter.POST("/publish/action/", controller.Publish) //这里是POST请求！
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", controller.FavoriteAction) //这里是POST请求！
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/action/", controller.CommentAction) //这里是POST请求！
	apiRouter.GET("/comment/list/", controller.CommentList)      //second

	// extra apis - II
	apiRouter.POST("/relation/action/", controller.RelationAction) //这里是POST请求！
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
