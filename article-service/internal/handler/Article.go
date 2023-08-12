package handler

import (
	"article-service/internal/service"
	"context"
)

type ArticleService struct {
}

func (a ArticleService) GetArticle(ctx context.Context, model *service.ArticleModel) (*service.ArticleDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

//func (a ArticleService) GetArticle(ctx context.Context, model *service.ArticleModel) (resq *service.ArticleDetailResponse, err error) {
//	Article, code := repository.GetArticle(model.Title)
//	resq.ArticleDetail = Article
//
//	return resq, nil
//
//}

func (a ArticleService) GetTypeArticles(ctx context.Context, model *service.ArticleModel) (*service.ArticleDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a ArticleService) GetRandArticle(ctx context.Context, model *service.ArticleModel) (*service.ArticleDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a ArticleService) DeleteArticle(ctx context.Context, model *service.ArticleModel) (*service.ArticleDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a ArticleService) GetArticleContent(ctx context.Context, model *service.ArticleModel) (*service.ArticleContentDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a ArticleService) DeleteArticleContent(ctx context.Context, model *service.ArticleModel) (*service.ArticleContentDetailResponse, error) {
	//TODO implement me
	panic("implement me")
}
