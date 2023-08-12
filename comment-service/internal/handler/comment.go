package handler

import (
	"comment-service/internal/repository"
	"comment-service/internal/service"
	"comment-service/pkg/errMsg"
	"context"
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
	var commentDate []*service.Comments
	comments, total := repository.GetComments(req.ArticleId, req.ArticleTitle)
	for _, i := range comments {
		data := repository.BuildComment(&i)
		commentDate = append(commentDate, data)
	}
	resq = new(service.CommentsDetailResponse)
	resq.Total = total
	resq.CommentDetail = commentDate
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
