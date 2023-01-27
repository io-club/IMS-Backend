package main

import (
	"context"
	demouser "fish_net/kitex_gen/demouser"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *demouser.CreateUserRequest) (resp *demouser.CreateUserResponse, err error) {
	// TODO: Your code here...
	return
}

// MGetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) MGetUser(ctx context.Context, req *demouser.MGetUserRequest) (resp *demouser.MGetUserResponse, err error) {
	// TODO: Your code here...
	return
}

// CheckUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *demouser.CheckUserRequest) (resp *demouser.CheckUserResponse, err error) {
	// TODO: Your code here...
	return
}
