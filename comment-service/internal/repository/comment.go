package repository

import (
	"comment-service/internal/service"
	"comment-service/pkg/errMsg"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Comments struct {
	gorm.Model
	Id        uint32
	UserId    uint32
	ArticleId string
	Title     string
	Username  string
	Content   string
	Img       string
	Agrees    int64
	Target    string
	Pid       int32
}

// Comment 评论

func AddComment(comment Comments) uint32 {
	err = db.Select("ID", "CreatedAt", "UpdatedAt", "DeletedAt", "Id", "UserId",
		"ArticleId", "Title", "Username", "Content", "Agrees", "Target", "Pid").Create(&comment).Error
	marshal, _ := json.Marshal(comment)
	err = redisCache.LPush(comment.ArticleId+comment.Title, marshal).Err()
	if err != nil {
		fmt.Println("数据库添加失败", err)
		return errMsg.ERROR
	}
	return errMsg.SUCCESS
}

func DeleteComment(articleId, title string) uint32 {
	var comment Comments
	// Unscoped不使用软删除
	err := db.Unscoped().Where("article_id= ? and title = ?", articleId, title).Delete(&comment).Error
	if err != nil {
		return errMsg.ERROR
	}

	time.Sleep(time.Millisecond * 500)
	err = redisCache.Del(articleId + title).Err()

	if err != nil {
		return errMsg.ERROR
	}
	return errMsg.SUCCESS
}

func GetComments(articleId, title string) ([]Comments, int64) {
	var comments []Comments
	var total int64
	result, err := redisCache.LRange(articleId+title, 0, 100).Result()
	if len(result) != 0 && err != nil {
		for _, i := range result {
			var comm Comments
			json.Unmarshal([]byte(i), &comm)
			comments = append(comments, comm)
		}
	} else {
		err = db.Raw("select * from comments where article_id=? and title =?", articleId, title).Scan(&comments).Error
		for _, i := range comments {
			marshal, _ := json.Marshal(i)
			err = redisCache.LPush(i.ArticleId+title, marshal).Err()
			if err != nil {
				log.Println("redis存储失败: ", err)
			}
		}
	}

	if err != nil {
		log.Println(err)
		return comments, 0
	}

	total = int64(len(comments))
	return comments, total
}

func AddAgree(id uint32, userId uint32, agrees int64) uint32 {
	err := db.Model(&Comments{}).Where("id=? and user_id =?", id, userId).Update("agrees", agrees).Error
	if err != nil {
		return errMsg.CommentAddWrong
	}
	return errMsg.SUCCESS
}
func BuildComment(comments *Comments) *service.Comments {
	return &service.Comments{
		Id:           comments.Id,
		UserId:       comments.UserId,
		ArticleId:    comments.ArticleId,
		ArticleTitle: comments.Title,
		Username:     comments.Username,
		Content:      comments.Content,
		Img:          comments.Img,
		Agrees:       comments.Agrees,
		Target:       comments.Target,
		Pid:          comments.Pid,
	}
}

func UnBuildComment(comments *service.Comments) *Comments {

	return &Comments{
		Id:        comments.Id,
		UserId:    comments.UserId,
		ArticleId: comments.ArticleId,
		Title:     comments.ArticleTitle,
		Username:  comments.Username,
		Content:   comments.Content,
		Img:       comments.Img,
		Agrees:    comments.Agrees,
		Target:    comments.Target,
		Pid:       comments.Pid,
	}
}
