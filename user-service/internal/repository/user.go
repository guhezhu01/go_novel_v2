package repository

import (
	"gorm.io/gorm"
	"user-service/internal/service"
	"user-service/pkg/errMsg"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null " json:"username" validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	Role     uint32 `gorm:"type:int;DEFAULT:2" json:"role" validate:"required,gte=2" label:"角色码"`
}

func (user *User) CheckUser(username string) (code uint32) {
	db.Select("id").Where("username", username).First(&user)
	if user.ID > 0 {
		return errMsg.UsernameUsed
	}
	return errMsg.SUCCESS
}

func (user *User) CreateUser(userData *User) (code uint32) {
	err = db.Create(&userData).Error
	if err != nil {
		return errMsg.ERROR //500
	}
	return errMsg.SUCCESS
}

func (user *User) GetUsers() ([]User, int) {
	var users []User
	var total int

	err = db.Model(&users).Find(&users).Error
	total = len(users)

	if err != nil {
		return users, 0
	}
	return users, total
}
func (user *User) GetUser(username string) (User, uint32) {
	err := db.Model(user).Where("username =?", username).Find(&user).Error
	if err != nil || user.Username == "" {
		return *user, errMsg.UserNotExist
	}
	return *user, errMsg.SUCCESS
}
func (user *User) DeleteUser(id int) uint32 {
	err = db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errMsg.ERROR
	}
	return errMsg.SUCCESS
}

func (user *User) EditUser(id int, userData *User) uint32 {
	var maps = make(map[string]interface{})
	maps["username"] = userData.Username
	maps["role"] = userData.Role
	err := db.Model(&user).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errMsg.ERROR
	}
	return errMsg.SUCCESS
}

func (user *User) CheckLogin(username string, password string) (User, uint32) {
	db.Where("username = ?", username).First(&user)
	if user.Username == "" {
		return *user, errMsg.UserNotExist
	} else if user.Password != password {
		return *user, errMsg.PasswordWrong
	}
	return *user, errMsg.SUCCESS
}

func (user *User) UpdateUserPassword(username, password string, id uint) uint32 {
	err := db.Model(&User{}).Where("username =? and id =?", username, id).Update("password", password).Error
	if err != nil {
		return errMsg.UpdatePasswordWrong
	}
	return errMsg.SUCCESS
}

func BuildUser(user User) *service.UserModel {
	userModel := service.UserModel{
		UserName: user.Username,
		Password: user.Password,
		Role:     user.Role,
	}
	return &userModel
}
