package service

import (
	"context"
	"errors"
	"time"

	"users-grpc/internal/models"
	"users-grpc/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserServiceServer struct {
	proto.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserServiceServer(db *gorm.DB) *UserServiceServer {
	return &UserServiceServer{db: db}
}

func toProto(u models.User) *proto.User {
	return &proto.User{
		Id:        u.ID,
		Name:      u.Name,
		Age:       u.Age,
		Address:   u.Address,
		CreatedAt: timestamppb.New(u.CreatedAt.UTC()),
		UpdatedAt: timestamppb.New(u.UpdatedAt.UTC()),
	}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, in *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	usr := models.User{Name: in.GetName(), Age: in.GetAge(), Address: in.GetAddress()}
	if err := s.db.Create(&usr).Error; err != nil {
		return nil, err
	}
	return &proto.CreateUserResponse{User: toProto(usr)}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, in *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	var usr models.User
	if err := s.db.First(&usr, in.GetId()).Error; err != nil {
		return nil, err
	}
	return &proto.GetUserResponse{User: toProto(usr)}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, in *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	var usr models.User
	if err := s.db.First(&usr, in.GetId()).Error; err != nil {
		return nil, err
	}
	if name := in.GetName(); name != "" {
		usr.Name = name
	}
	if in.Age != 0 {
		usr.Age = in.Age
	}
	if addr := in.GetAddress(); addr != "" {
		usr.Address = addr
	}
	usr.UpdatedAt = time.Now().UTC()
	if err := s.db.Save(&usr).Error; err != nil {
		return nil, err
	}
	return &proto.UpdateUserResponse{User: toProto(usr)}, nil
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, in *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	res := s.db.Delete(&models.User{}, in.GetId())
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	return &proto.DeleteUserResponse{Ok: true}, nil
}

func (s *UserServiceServer) ListUsers(ctx context.Context, in *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	page := in.GetPage()
	if page < 1 {
		page = 1
	}
	size := in.GetPageSize()
	if size <= 0 || size > 100 {
		size = 10
	}

	var users []models.User
	var total int64
	s.db.Model(&models.User{}).Count(&total)

	if err := s.db.Order("id ASC").Limit(int(size)).Offset(int((page - 1) * size)).Find(&users).Error; err != nil {
		return nil, err
	}

	out := make([]*proto.User, 0, len(users))
	for _, u := range users {
		out = append(out, toProto(u))
	}

	return &proto.ListUsersResponse{Users: out, Total: total}, nil
}
