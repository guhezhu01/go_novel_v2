package handler

import (
	"context"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/errMsg"
)

type UserService struct {
}

func (s *UserService) GetUser(context.Context, *service.UserRequest) (*service.UserDetailResponse, error) {

	panic("implement me")
}

func (s *UserService) DeleteUser(context.Context, *service.UserRequest) (*service.UserDetailResponse, error) {

	panic("implement me")
}

func (s *UserService) EditUser(context.Context, *service.UserRequest) (*service.UserDetailResponse, error) {

	panic("implement me")
}

func (s *UserService) UpdateUserPassword(context.Context, *service.UserRequest) (*service.UserDetailResponse, error) {

	panic("implement me")
}

func (s *UserService) GetUsers(context.Context, *service.UserRequest) (*service.UserDetailResponse, error) {

	panic("implement me")
}

//func NewUserService() *UserService {
//	return &UserService{}
//}

func (*UserService) UserLogin(_ context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {
	var user repository.User
	var code uint32
	resq = new(service.UserDetailResponse)
	code = user.CheckUser(req.GetUserDetail().GetUserName())
	if code == errMsg.SUCCESS {
		resq.Code = errMsg.UserNotExist
		resq.Msg = errMsg.GetErrMsg(errMsg.UserNotExist)
		return resq, nil
	}
	user, code = user.CheckLogin(req.GetUserDetail().GetUserName(), req.GetUserDetail().GetPassword())
	resq.UserDetail = repository.BuildUser(user)
	resq.Msg = errMsg.GetErrMsg(code)
	return resq, nil
}

func (*UserService) UserRegister(_ context.Context, req *service.UserRequest) (resq *service.UserDetailResponse, err error) {
	var user repository.User
	var code uint32
	user.Username = req.GetUserDetail().GetUserName()
	user.Password = req.GetUserDetail().GetPassword()
	user.Role = req.GetUserDetail().GetRole()
	resq = new(service.UserDetailResponse)
	code = user.CheckUser(req.GetUserDetail().GetUserName())
	if code == 200 {
		code = user.CreateUser(&user)
	}
	resq.Code = code
	resq.UserDetail = repository.BuildUser(user)
	resq.Msg = errMsg.GetErrMsg(code)
	return resq, nil
}
