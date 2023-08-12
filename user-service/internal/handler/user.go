package handler

import (
	"context"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/errMsg"
)

type UserService struct {
}

// GetUser 获取用户
func (s *UserService) GetUser(ctx context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {

	var user repository.User
	var code uint32
	resq = new(service.UserDetailResponse)
	user, code = user.GetUser(ctx, req.GetUserDetail().Username)
	resq.UserDetail = repository.BuildUser(user)
	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)
	return resq, nil

}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {
	var user repository.User
	var code uint32
	resq = new(service.UserDetailResponse)
	code = user.DeleteUser(ctx, req.UserDetail.UserId)
	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)
	return resq, nil

}

// EditUser 修改用户
func (s *UserService) EditUser(ctx context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {
	var user repository.User
	var code uint32
	resq = new(service.UserDetailResponse)
	code = user.CheckUser(ctx, req.UserDetail.Username)
	if code != errMsg.SUCCESS {
		code = user.EditUser(ctx, int(req.GetUserDetail().UserId), req.UserDetail)
		resq.UserDetail = repository.BuildUser(user)
	}

	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)
	return resq, nil
}

// UpdateUserPassword 修改用户密码
func (s *UserService) UpdateUserPassword(ctx context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {
	var user repository.User
	var code uint32
	resq = new(service.UserDetailResponse)
	code = user.UpdateUserPassword(ctx, req.GetUserDetail().Username, req.GetUserDetail().Password, uint(req.GetUserDetail().UserId))
	resq.UserDetail = repository.BuildUser(user)
	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)
	return resq, nil
}

// GetUsers 获取所有用户
func (s *UserService) GetUsers(context.Context, *service.UserRequest) (*service.UserDetailResponse, error) {
	panic("implement me")
}

// UserLogin 用户登录验证
func (s *UserService) UserLogin(ctx context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {
	var user repository.User
	var code uint32
	resq = new(service.UserDetailResponse)
	code = user.CheckUser(ctx, req.GetUserDetail().GetUsername())
	if code == errMsg.SUCCESS {
		resq.Code = errMsg.UserNotExist
		resq.Msg = errMsg.GetErrMsg(errMsg.UserNotExist)
		return resq, nil
	}
	user, code = user.CheckLogin(ctx, req.GetUserDetail().GetUsername(), req.GetUserDetail().GetPassword())

	resq.UserDetail = repository.BuildUser(user)
	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)

	return resq, nil
}

// UserRegister 注册用户
func (*UserService) UserRegister(ctx context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {
	var user repository.User
	var code uint32
	user.Username = req.GetUserDetail().GetUsername()
	user.Password = req.GetUserDetail().GetPassword()
	user.Role = req.GetUserDetail().GetRole()
	resq = new(service.UserDetailResponse)
	code = user.CheckUser(ctx, req.GetUserDetail().GetUsername())
	if code == 200 {
		code = user.CreateUser(ctx, &user)
	}
	resq.UserDetail = repository.BuildUser(user)
	resq.Code = code
	resq.Msg = errMsg.GetErrMsg(code)
	return resq, nil
}
