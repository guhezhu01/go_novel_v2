package repository

import (
	"context"
	"gorm.io/gorm"
	"user-service/internal/service"
	"user-service/pkg/errMsg"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null "  validate:"required,min=4,max=12" label:"用户名"`
	Password string `gorm:"type:varchar(500);not null"  validate:"required,min=6,max=120" label:"密码"`
	Role     uint32 `gorm:"type:int;DEFAULT:2" validate:"required,gte=2" label:"角色码"`
}

// CheckUser 检查用户是否存在，未注册返回200
func (user *User) CheckUser(ctx context.Context, username string) (code uint32) {
	db.WithContext(ctx).Select("id").Where("username", username).First(&user)
	if user.ID > 0 {
		return errMsg.UsernameUsed
	}
	return errMsg.SUCCESS
}

// CreateUser 创建用户
func (user *User) CreateUser(ctx context.Context, userData *User) (code uint32) {
	err = db.WithContext(ctx).Create(&userData).Error
	if err != nil {
		return errMsg.ERROR
	}
	return errMsg.SUCCESS
}

// GetUsers 获取所有用户
func (user *User) GetUsers(ctx context.Context) ([]User, int) {
	var users []User
	var total int

	err = db.WithContext(ctx).Model(&users).Find(&users).Error
	total = len(users)

	if err != nil {
		return users, 0
	}
	return users, total
}

// GetUser 获取用户
func (user *User) GetUser(ctx context.Context, username string) (User, uint32) {
	err := db.WithContext(ctx).Model(user).Where("username =?", username).Find(&user).Error
	if err != nil || user.Username == "" {
		return *user, errMsg.UserNotExist
	}
	return *user, errMsg.SUCCESS
}

// DeleteUser 删除用户
func (user *User) DeleteUser(ctx context.Context, id uint32) uint32 {
	err = db.WithContext(ctx).Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errMsg.ERROR
	}
	return errMsg.SUCCESS
}

// EditUser 修改用户
func (user *User) EditUser(ctx context.Context, id int, userData *service.UserModel) uint32 {
	var maps = make(map[string]interface{})
	maps["username"] = userData.Username
	maps["role"] = userData.Role
	err := db.WithContext(ctx).Model(&user).Where("id = ? ", id).Updates(maps).Error
	if err != nil {
		return errMsg.ERROR
	}
	return errMsg.SUCCESS
}

// CheckLogin 检查用户密码
func (user *User) CheckLogin(ctx context.Context, username string, password string) (User, uint32) {

	db.WithContext(ctx).Where("username = ?", username).First(&user)
	if user.Username == "" {
		return *user, errMsg.UserNotExist
	} else if user.Password != password {
		return *user, errMsg.PasswordWrong
	}
	return *user, errMsg.SUCCESS
}

// UpdateUserPassword 更新密码
func (user *User) UpdateUserPassword(ctx context.Context, username, password string, id uint) uint32 {
	err := db.WithContext(ctx).Model(&User{}).Where("username =? and id =?", username, id).Update("password", password).Error
	if err != nil {
		return errMsg.UpdatePasswordWrong
	}
	return errMsg.SUCCESS
}

func BuildUser(user User) *service.UserModel {
	userModel := service.UserModel{
		Username: user.Username,
		Password: user.Password,
		Role:     user.Role,
		UserId:   uint32(user.ID),
	}
	return &userModel
}
