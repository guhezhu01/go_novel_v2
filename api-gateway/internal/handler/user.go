package handler

import (
	"api-gateway/discovery"
	"api-gateway/internal/service"
	"api-gateway/middleware"
	"api-gateway/pkg/errMsg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var userRes service.UserRequest
	_ = c.ShouldBindJSON(&userRes.UserDetail)
	fmt.Println(c.Request.Context().Value("to"))
	//c.Request = c.Request.WithContext(ctx)
	getService := discovery.GetService(viper.GetString("UserService.name"), viper.GetString("UserService.tag"))
	grpcClient := getService.(service.UserServiceClient)
	//调用服务端的方法

	p, _ := grpcClient.UserLogin(c.Request.Context(), &userRes)
	if p.Code == errMsg.SUCCESS {
		userRes.GetUserDetail().Token, _ = middleware.SetToken(userRes.GetUserDetail().Username)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          p.Code,
		"message":         errMsg.GetErrMsg(p.Code),
		"userInformation": userRes.UserDetail,
	})
}

//// GetUsers 获取所有用户
//func GetUsers(c *gin.Context) {
//
//	code = errMsg.SUCCESS
//	c.JSON(
//		http.StatusOK, gin.H{
//			"status":  code,
//			"data":    data,
//			"total":   total,
//			"message": errMsg.GetErrMsg(code),
//		},
//	)
//}
//

// GetUser 查询单个用户
func GetUser(c *gin.Context) {
	var userRes service.UserRequest

	var user service.UserModel
	user.Username = c.Query("username")
	userRes.UserDetail = &user
	getService := discovery.GetService(viper.GetString("UserService.name"), viper.GetString("UserService.tag"))
	grpcClient := getService.(service.UserServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.GetUser(c.Request.Context(), &userRes)

	c.JSON(http.StatusOK, gin.H{
		"status":  p.Code,
		"data":    p.UserDetail,
		"message": errMsg.GetErrMsg(p.Code),
	})
}

// UserRegister 添加用户
func UserRegister(c *gin.Context) {
	var userRes service.UserRequest
	_ = c.ShouldBindJSON(userRes.UserDetail)
	getService := discovery.GetService(viper.GetString("UserService.name"), viper.GetString("UserService.tag"))
	grpcClient := getService.(service.UserServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.UserRegister(c.Request.Context(), &userRes)
	if p.Code != errMsg.SUCCESS {
		p.UserDetail = nil
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  p.Code,
		"data":    p.UserDetail,
		"message": errMsg.GetErrMsg(p.Code),
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	var userRes service.UserRequest
	var user service.UserModel
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	user.UserId = uint32(id)
	userRes.UserDetail = &user
	getService := discovery.GetService(viper.GetString("UserService.name"), viper.GetString("UserService.tag"))
	grpcClient := getService.(service.UserServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.DeleteUser(c.Request.Context(), &userRes)
	c.JSON(http.StatusOK, gin.H{
		"status":  p.Code,
		"message": errMsg.GetErrMsg(p.Code),
	})
}

// EditUser 修改用户
func EditUser(c *gin.Context) {
	var userRes service.UserRequest
	_ = c.ShouldBindJSON(&userRes.UserDetail)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userRes.UserDetail.UserId = uint32(id)
	getService := discovery.GetService(viper.GetString("UserService.name"), viper.GetString("UserService.tag"))
	grpcClient := getService.(service.UserServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.EditUser(c.Request.Context(), &userRes)
	c.JSON(http.StatusOK, gin.H{
		"status":  p.Code,
		"message": errMsg.GetErrMsg(p.Code)})
}
