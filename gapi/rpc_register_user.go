package gapi

import (
	"context"
	"encoding/base64"
	"time"

	pb_reg "github.com/Streamfair/common_proto/IdentityProvider/pb/register"
	user_pb "github.com/Streamfair/common_proto/UserService/pb"
	user "github.com/Streamfair/common_proto/UserService/pb/user"
	"github.com/Streamfair/streamfair_idp/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	IDP_svc_address  = "streamfair_idp:9091"
	USER_svc_address = "streamfair_user_service:9094"
)

// RegisterUser registers a user in the usersdb and returns the user.

// Microservices involved: UserService
// 1. Creates the user with the given information via UserService.
// 5. Return the user information.
func (server *Server) RegisterUser(ctx context.Context, req *pb_reg.RegisterUserRequest) (*pb_reg.RegisterUserResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               IDP_svc_address,
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	}

	pool := NewClientPool(poolConfig)

	user, err := registerUser(ctx, pool, USER_svc_address, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}

	rps := &pb_reg.RegisterUserResponse{
		User: ConvertRegisteredUser(user),
	}

	return rps, nil
}

func registerUser(ctx context.Context, pool *ConnectionPool, address string, req *pb_reg.RegisterUserRequest) (*user.User, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}

	client := user_pb.NewUserServiceClient(conn)

	byteHash, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to prepare registering %v", err)
	}
	hashedPassword := base64.StdEncoding.EncodeToString(byteHash.Hash)
	passwordSalt := base64.StdEncoding.EncodeToString(byteHash.Salt)

	arg := &user.CreateUserRequest{
		Username:     req.GetUsername(),
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		PasswordSalt: passwordSalt,
		CountryCode:  req.GetCountryCode(),
		RoleId:       int64(req.GetRoleId()),
		Status:       req.GetStatus(),
	}
	resp, err := client.CreateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return resp.User, nil
}
