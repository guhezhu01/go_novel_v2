package routes

import (
	"api-gateway/middleware"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Cors(), middleware.Logger())
	auth := engine.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{
		//auth.DELETE("user/delete/:id", v1.DeleteUser)
		//auth.PUT("user/update/:id", v1.EditUser)
		//auth.GET("users", v1.GetUsers)
		//auth.GET("user", v1.GetUser)
		//auth.POST("user/update/password", v1.UpdateUserPassword)
		//
		//userInfor := auth.Group("user/information")
		//{
		//	userInfor.GET("/", v1.GetUserInformation)
		//	userInfor.POST("add", v1.CreateUserInformation)
		//	userInfor.PUT("update", v1.EditUserInformation)
		//}
		//
		//userColl := auth.Group("user/collection")
		//{
		//	userColl.GET("/", v1.GetUserCollections)
		//	userColl.POST("add", v1.AddUserCollection)
		//	userColl.DELETE("delete", v1.DeleteUserCollection)
		//}
	}
	//
	//router := engine.Group("api/v1")
	//{
	//	router.GET("article/download", v1.DownloadArticle)
	//	router.DELETE("article/delete", v1.DeleteArticle)
	//
	//	router.POST("user/add", v1.AddUser)
	//	router.POST("login", v1.Login)
	//
	//	router.POST("article", v1.GetArticle)
	//	router.GET("article/type", v1.GetTypeArticles)
	//	router.GET("article/rand", v1.GetRandArticle)
	//
	//	router.GET("comments", v1.GetComments)
	//	router.DELETE("comment/delete/:id/:user_id", v1.DeleteComment)
	//	router.POST("comment/add", v1.AddComment)
	//
	//	router.POST("comment/add/agrees/:id/:user_id", v1.AddAgree)
	//}
	return engine
}
