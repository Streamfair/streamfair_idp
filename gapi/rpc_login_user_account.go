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
	"github.com/Streamfair/streamfair_idp/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) LoginUserAccount(ctx context.Context, req *pb_login.LoginUserAccountRequest) (*pb_login.LoginUserAccountResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               IdpSvcAddress,
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	}

	pool := NewClientPool(poolConfig)

	username := req.GetUsername()

	// Fetch user from the database
	userAccount, err := server.store.GetUserAccountByUserAccountname(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user_account not found: %v", err)
		}
		return nil, handleDatabaseError(err)
	}

	byteHash, err := base64.StdEncoding.DecodeString(userAccount.PasswordHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "decoding error occured.")
	}
	byteSalt, err := base64.StdEncoding.DecodeString(userAccount.PasswordSalt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "decoding error occured.")
	}

	err = util.ComparePassword(byteHash, byteSalt, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password.")
	}

	accessToken, err := createToken(ctx, pool, TOKEN_svc_address, &token.CreateTokenRequest{
		UserId:    userAccount.ID,
		ExpiresAt: server.config.AccessTokenDuration.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token: %v.", err)
	}

	refreshToken, err := createRefreshToken(ctx, pool, TOKEN_svc_address, &refreshToken.CreateRefreshTokenRequest{
		UserId:    userAccount.ID,
		ExpiresAt: server.config.RefreshTokenDuration.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token: %v.", err)
	}

	mtdt := server.extractMetadata(ctx)
	args := &session.CreateSessionRequest{
		Uuid:         refreshToken.Payload.GetUuid(),
		Username:     userAccount.Username,
		RefreshToken: refreshToken.RefreshToken.GetToken(),
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshToken.Payload.GetExpiredAt(),
	}
	session, err := createSession(ctx, pool, SESSION_svc_address, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v.", err)
	}

	rps := &pb_login.LoginUserAccountResponse{
		UserAccount: ConvertUserAccount(userAccount),
		Session: &pb_login.SessionData{
			SessionId:             session.GetUuid(),
			IsBlocked:             session.GetIsBlocked(),
			UserAgent:             session.GetUserAgent(),
			ClientIp:              session.GetClientIp(),
			AccessToken:           accessToken.GetToken().GetToken(),
			AccessTokenExpiresAt:  accessToken.GetToken().GetExpiresAt(),
			RefreshToken:          refreshToken.GetRefreshToken().GetToken(),
			RefreshTokenExpiresAt: refreshToken.GetPayload().GetExpiredAt(),
			SessionCreatedAt:      session.GetCreatedAt(),
			SessionExpiresAt:      session.GetExpiresAt(),
		},
	}

	return rps, nil
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
