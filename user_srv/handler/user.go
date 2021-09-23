package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"mxshop/user_srv/global"
	"mxshop/user_srv/model"
	"mxshop/user_srv/proto"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"

	//"google.golang.org/protobuf/types/known/emptypb"

	"github.com/anaskhan96/go-password-encoder"

	"google.golang.org/grpc/codes"
	//"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/status"

	"gorm.io/gorm"
)

type UserServer struct{}

func ModelToRsponse(user model.User) proto.UserInfoResponse {
	// 在grpc的message中字段有默认值，你不能随便赋值nil进去， 容易报错
	// 这里要搞清楚，哪些字段是有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil { // 如果有值，就覆盖， 没有有为默认0
		userInfoRsp.BirthDay = uint32(user.Birthday.Unix())
	}
	return userInfoRsp
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// 获取用户列表
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected) // 全部数据

	// 分页
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users) // 分页的数据
	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp)
	}

	return rsp, nil
}

// 通过手机号码查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

// 通过ID查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

// 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	// 1. 检查用户是否存在
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	// 密码加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	user.Password = newPassword

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

// 更新用户
func (s UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0) // int64 转换成时间格式
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

// 校验密码
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
