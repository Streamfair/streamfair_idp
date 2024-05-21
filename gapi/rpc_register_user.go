package gapi

import (
	"context"
	"time"

	pb "github.com/Streamfair/common_proto/IdentityProvider/pb/register"
	user_pb "github.com/Streamfair/common_proto/UserService/pb"
	user "github.com/Streamfair/common_proto/UserService/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	IDP_svc_address = "streamfair_idp:9091"
	USER_svc_address = "streamfair_user_service:9094"
)

// RegisterUser registers a user in the usersdb and returns the user.

// Microservices involved: UserService
// 1. Creates the user with the given information via UserService.
// 5. Return the user information.
func (server *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               IDP_svc_address,
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	}


	pool := NewClientPool(poolConfig)

	user, err := createUser(ctx, pool, USER_svc_address, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}

	rps := &pb.RegisterUserResponse{
		User: user,
	}

	return rps, nil
}

func createUser(ctx context.Context, pool *ConnectionPool, address string, params *pb.RegisterUserRequest) (*user.User, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}

	client := user_pb.NewUserServiceClient(conn)

	

	req := &user.CreateUserRequest{
		Username: params.GetUsername(),
		FullName: params.GetFullName(),
		Email: params.GetEmail(),
		Password: params.GetPassword(),
		RoleId: int64(params.GetRoleId()),
		Status: params.GetStatus(),
	}
	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return resp.User, nil
}
