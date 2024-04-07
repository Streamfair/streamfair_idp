package gapi

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	pb "github.com/Streamfair/common_proto/IdentityProvider/pb/login"
	session_pb "github.com/Streamfair/common_proto/SessionService/pb"
	session "github.com/Streamfair/common_proto/SessionService/pb/session"
	token_pb "github.com/Streamfair/common_proto/TokenService/pb"
	refreshToken "github.com/Streamfair/common_proto/TokenService/pb/refresh_token"
	token "github.com/Streamfair/common_proto/TokenService/pb/token"
	user_pb "github.com/Streamfair/common_proto/UserService/pb"
	user "github.com/Streamfair/common_proto/UserService/pb/user"
	"github.com/Streamfair/streamfair_idp/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getUser(ctx context.Context, pool *ConnectionPool, address string, username string) (*user.User, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}

	client := user_pb.NewUserServiceClient(conn)

	req := &user.GetUserByValueRequest{
		Username: username,
	}
	resp, err := client.GetUserByValue(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return resp.User, nil
}

func createSession(ctx context.Context, pool *ConnectionPool, address string, args *session.CreateSessionRequest) (*session.Session, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to SessionService: %v", err)
	}

	client := session_pb.NewSessionServiceClient(conn)

	resp, err := client.CreateSession(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return resp.Session, nil
}

func createToken(ctx context.Context, pool *ConnectionPool, address string, args *token.CreateTokenRequest) (*token.CreateTokenResponse, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to SessionService: %v", err)
	}

	client := token_pb.NewTokenServiceClient(conn)

	resp, err := client.CreateToken(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return resp, nil
}

func createRefreshToken(ctx context.Context, pool *ConnectionPool, address string, args *refreshToken.CreateRefreshTokenRequest) (*refreshToken.CreateRefreshTokenResponse, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to SessionService: %v", err)
	}

	client := token_pb.NewTokenServiceClient(conn)

	resp, err := client.CreateRefreshToken(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return resp, nil
}

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               "your_user_service_address",
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	}
	pool := NewClientPool(poolConfig)

	user, err := getUser(ctx, pool, "https://streamfair_user_service:9094/streamfair/v1/get_user_by_value", req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "user not found: %v", err)
	}

	byteHash, err := base64.StdEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "decoding error occured.")
	}
	byteSalt, err := base64.StdEncoding.DecodeString(user.PasswordSalt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "decoding error occured.")
	}

	err = util.ComparePassword(byteHash, byteSalt, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password.")
	}

	accessToken, err := createToken(ctx, pool, "https://streamfair_token_service:9092/streamfair/v1/create_token", &token.CreateTokenRequest{
		UserId:    user.Id,
		ExpiresAt: server.config.AccessTokenDuration.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v.", err)
	}

	refreshToken, err := createRefreshToken(ctx, pool, "https://streamfair_token_service:9092/streamfair/v1/create_token", &refreshToken.CreateRefreshTokenRequest{
		UserId:    user.Id,
		ExpiresAt: server.config.RefreshTokenDuration.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v.", err)
	}

	mtdt := server.extractMetadata(ctx)
	args := &session.CreateSessionRequest{
		Uuid:         refreshToken.Payload.Uuid,
		Username:     user.Username,
		RefreshToken: refreshToken.RefreshToken.Token,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshToken.Payload.ExpiredAt,
	}
	session, err := createSession(ctx, pool, "https://streamfair_session_service:9093/streamfair/v1/create_session", args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v.", err)
	}

	rps := &pb.LoginUserResponse{
		User:                  user,
		SessionId:             session.Uuid,
		AccessToken:           accessToken.Token.Token,
		RefreshToken:          refreshToken.RefreshToken.Token,
		AccessTokenExpiresAt:  accessToken.Payload.ExpiredAt,
		RefreshTokenExpiresAt: refreshToken.Payload.ExpiredAt,
	}

	return rps, nil
}
