package handler

import (
	"api-gateway/discovery"
	"api-gateway/internal/service"
	"api-gateway/pkg/errMsg"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func AddComment(c *gin.Context) {
	var comment service.Comments
	_ = c.ShouldBindJSON(&comment)
	getService := discovery.GetService(c.Request.Context(), "comment service", "grpc")
	grpcClient := getService.(service.CommentServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.AddComment(context.TODO(), &comment)
	c.JSON(http.StatusOK, gin.H{
		"status":  p.Code,
		"message": errMsg.GetErrMsg(p.Code),
	})

}

func DeleteComment(c *gin.Context) {
	var comment service.Comments
	_ = c.ShouldBindJSON(&comment)
	getService := discovery.GetService(c.Request.Context(), "comment-service", "grpc")
	grpcClient := getService.(service.CommentServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.DeleteComment(context.TODO(), &comment)
	c.JSON(http.StatusOK, gin.H{
		"status":  p.Code,
		"message": errMsg.GetErrMsg(p.Code),
	})
}

func GetComments(c *gin.Context) {
	var comment service.Comments
	comment.ArticleId = c.Query("article_id")
	comment.ArticleTitle = c.Query("title")

	getService := discovery.GetService(c.Request.Context(), "comment-service", "grpc")
	grpcClient := getService.(service.CommentServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.GetComments(c, &comment)
	if p.Total > 0 {
		p.Code = errMsg.SUCCESS
	} else {
		p.Code = errMsg.CommentNotExist
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   p.Code,
		"message":  errMsg.GetErrMsg(p.Code),
		"comments": p.CommentDetail,
		"total":    p.Total,
	})

}

func AddAgree(c *gin.Context) {
	var comment service.Comments
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	comment.Id = uint32(id)
	userId, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
	comment.UserId = uint32(userId)
	_ = c.ShouldBindJSON(&comment)
	getService := discovery.GetService(c.Request.Context(), "comment service", "grpc")
	grpcClient := getService.(service.CommentServiceClient)
	//调用服务端的方法
	p, _ := grpcClient.AddAgree(context.TODO(), &comment)
	c.JSON(http.StatusOK, gin.H{
		"status":  p.Code,
		"message": errMsg.GetErrMsg(p.Code),
	})
}
