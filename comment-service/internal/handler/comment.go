package handler

import (
	"comment-service/init-config"
	"comment-service/internal/repository"
	"comment-service/internal/service"
	"comment-service/middleware"
	"comment-service/pkg/errMsg"
	"context"
	"github.com/spf13/viper"
	"time"
)

type CommentsService struct {
}

func (c CommentsService) AddComment(ctx context.Context, req *service.Comments) (resq *service.CommentsDetailResponse, err error) {
	comment := repository.UnBuildComment(req)
	code := repository.AddComment(*comment)
	resq = new(service.CommentsDetailResponse)
	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)
	ctx.Done()
	return resq, nil
}

func (c CommentsService) DeleteComment(ctx context.Context, req *service.Comments) (resq *service.CommentsDetailResponse, err error) {
	code := repository.DeleteComment(req.ArticleId, req.ArticleTitle)
	resq = new(service.CommentsDetailResponse)
	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)
	ctx.Done()
	return resq, nil

}

func (c CommentsService) GetComments(ctx context.Context, req *service.Comments) (resq *service.CommentsDetailResponse, err error) {
	// 将[]service.CommentModel转为[]*service.CommentModel
	id, _ := middleware.GetRequestID(ctx, viper.GetString("consul.Name"))
	initConfig.DistributedLockConn.TryLock(viper.GetString("consul.Name"), time.Second, id)
	defer func() {
		initConfig.DistributedLockConn.Unlock(id)
	}()
	var commentDate []*service.Comments
	comments, total := repository.GetComments(req.ArticleId, req.ArticleTitle)
	for _, i := range comments {
		data := repository.BuildComment(&i)
		commentDate = append(commentDate, data)
	}
	resq = new(service.CommentsDetailResponse)
	resq.Total = total
	resq.CommentDetail = commentDate
	resq.Code = errMsg.SUCCESS
	ctx.Done()
	return resq, nil
}

func (c CommentsService) AddAgree(ctx context.Context, req *service.Comments) (resq *service.CommentsDetailResponse, err error) {
	code := repository.AddAgree(req.Id, req.UserId, req.Agrees)
	resq = new(service.CommentsDetailResponse)
	resq.Code = code
	ctx.Done()
	return resq, nil
}
