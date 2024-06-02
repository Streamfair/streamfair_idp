package gapi

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	pb_login "github.com/Streamfair/common_proto/IdentityProvider/pb/login"
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

const (
	TOKEN_svc_address   = "streamfair_token_service:9092"
	SESSION_svc_address = "streamfair_session_service:9093"
)

// LoginUser authenticates a user and returns a session token.

// Microservices involved: UserService, TokenService, SessionService
// 1. Get user by username from UserService.
// 2. Compare the password hash with the password provided.
// 3. Create an access token and a refresh token.
// 4. Create a session with the refresh token.
// 5. Return the user, session ID, access token, refresh token, and their expiration times.
func (server *Server) LoginUser(ctx context.Context, req *pb_login.LoginUserRequest) (*pb_login.LoginUserResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               IDP_svc_address,
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	}

	pool := NewClientPool(poolConfig)

	username := req.GetUsername()

	user, err := getUser(ctx, pool, USER_svc_address, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
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

	accessToken, err := createToken(ctx, pool, TOKEN_svc_address, &token.CreateTokenRequest{
		UserId:    user.Id,
		ExpiresAt: server.config.AccessTokenDuration.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v.", err)
	}

	refreshToken, err := createRefreshToken(ctx, pool, TOKEN_svc_address, &refreshToken.CreateRefreshTokenRequest{
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
	session, err := createSession(ctx, pool, SESSION_svc_address, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v.", err)
	}

	rps := &pb_login.LoginUserResponse{
		User:                  ConvertLoggedInUser(user),
		SessionId:             session.Uuid,
		AccessToken:           accessToken.Token.Token,
		RefreshToken:          refreshToken.RefreshToken.Token,
		AccessTokenExpiresAt:  accessToken.Payload.ExpiredAt,
		RefreshTokenExpiresAt: refreshToken.Payload.ExpiredAt,
	}

	return rps, nil
}

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
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
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
