package repository

import (
	"article-service/internal/service"
	"article-service/pkg/errMsg"
)

func GetTypeArticles(type_ string) ([]service.ArticleModel, int, int) {
	var articles []service.ArticleModel
	var total int

	err = db.Model(&articles).Where("type LIKE ?", type_+"%").Find(&articles).Error
	total = len(articles)

	if err != nil {
		return articles, 0, errMsg.CateNotExist
	}
	return articles, total, errMsg.SUCCESS
}

func GetArticle(title string) (service.ArticleModel, uint32) {
	var article service.ArticleModel
	err := db.Where("title = ?", title).First(&article).Error
	if err != nil {
		return article, errMsg.ArtNotExist
	}
	return article, errMsg.SUCCESS
}

//func GetRandArticle() ([]Article, int, int) {
//	var data []Article
//	var total int
//	sql := "SELECT * FROM article WHERE 1  ORDER BY rand() limit 12;"
//	err := db.Raw(sql).Scan(&data).Error
//	total = len(data)
//	if err != nil {
//		return data, 0, errMsg.ERROR
//	}
//	return data, total, errMsg.SUCCESS
//}
//
//func DeleteArticle(id string, type_ string) int {
//	var article Article
//	err := db.Unscoped().Where("id =? and type=?", id, type_).Delete(&article).Error
//	if err != nil {
//		return errMsg.ERROR
//	}
//	return errMsg.SUCCESS
//}
