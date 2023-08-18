package routes

import (
	"api-gateway/internal/handler"
	"api-gateway/middleware"
	"api-gateway/pkg/errMsg"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(1*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(testResponse),
	)
}

func testResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": errMsg.TIMEOUT,
		"msg":  errMsg.GetErrMsg(errMsg.TIMEOUT),
	})
}

func InitRoutes() *gin.Engine {

	engine := gin.Default()
	//middleware.Limit(engine)
	engine.Use(middleware.Cors(), middleware.Logger())
	auth := engine.Group("api/v1")
	auth.Use(middleware.JwtToken())
	{

		auth.DELETE("user/delete/:id", handler.DeleteUser)
		auth.PUT("user/update/:id", handler.EditUser)
		//auth.GET("users", handler.GetUsers)
		auth.GET("user", handler.GetUser)

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
	router := engine.Group("api/v1")
	{
		//	router.GET("article/download", v1.DownloadArticle)
		//	router.DELETE("article/delete", v1.DeleteArticle)
		//
		router.POST("user/add", handler.UserRegister)
		router.POST("login", handler.UserLogin)
		//
		//	router.POST("article", v1.GetArticle)
		//	router.GET("article/type", v1.GetTypeArticles)
		//	router.GET("article/rand", v1.GetRandArticle)
		//
		router.GET("comments", handler.GetComments)
		router.DELETE("comment/delete/", handler.DeleteComment)
		router.POST("comment/add", handler.AddComment)
		router.POST("comment/add/agrees/:id/:user_id", handler.AddAgree)
	}
	return engine
}
