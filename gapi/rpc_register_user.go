package gapi

import (
	"context"
	"time"

	"github.com/Streamfair/common_proto/IdentityProvider/pb/register";
	user_pb "github.com/Streamfair/common_proto/UserService/pb"
	user "github.com/Streamfair/common_proto/UserService/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoginUser authenticates a user and returns a session token.

// Microservices involved: UserService, TokenService, SessionService
// 1. Get user by username from UserService.
// 2. Compare the password hash with the password provided.
// 3. Create an access token and a refresh token.
// 4. Create a session with the refresh token.
// 5. Return the user, session ID, access token, refresh token, and their expiration times.
func (server *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserResponse) (*pb.RegisterUserResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               "streamfair_idp:9091",
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	}

	pool := NewClientPool(poolConfig)

	username := req.GetUsername()

	user, err := createUser(ctx, pool, "streamfair_user_service:9094", username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}

	rps := &pb.RegisterUserResponse{
		User: user,
	}

	return rps, nil
}

func createUser(ctx context.Context, pool *ConnectionPool, address string, username string) (*user.User, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}

	client := user_pb.NewUserServiceClient(conn)

	req := &user.CreateUserRequest{
		Username: username,
		// implement more fields
	}
	resp, err := client.CreateUser(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return resp.User, nil
}
