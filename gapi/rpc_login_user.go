package gapi

import (
	"context"
	"database/sql"
	"encoding/base64"

	pb "github.com/Streamfair/streamfair_idp/pb"
	pb_login "github.com/Streamfair/streamfair_idp/pb/login"
	"github.com/Streamfair/streamfair_idp/util"
	pb_session "github.com/Streamfair/streamfair_session_svc/pb/session"
	pb_user "github.com/Streamfair/streamfair_user_svc/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func getUser(ctx context.Context, pool *ConnectionPool, address string, username string) (*pb_user.User, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}

	client := pb.NewUserServiceClient(conn)

	req := &pb_user.GetUserByValueRequest{
		Username: username,
	}
	resp, err := client.GetUserByValue(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return resp.User, nil
}

func createSession(ctx context.Context, pool *ConnectionPool, address string, args *pb_session.CreateSessionRequest) (*pb_session.Session, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to SessionService: %v", err)
	}

	client := pb.NewSessionServiceClient(conn)

	resp, err := client.CreateSession(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return resp.Session, nil
}

func (server *Server) LoginUser(ctx context.Context, req *pb_login.LoginUserRequest) (*pb_login.LoginUserResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               "your_user_service_address",
		ConfigOptions:         []grpc.DialOption{},
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

	accessToken, accessPayload, err := server.localTokenMaker.CreateLocalToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token.")
	}

	refreshToken, refreshPayload, err := server.localTokenMaker.CreateLocalToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token.")
	}

	mtdt := server.extractMetadata(ctx)
	args := &pb_session.CreateSessionRequest{
		Uuid:         refreshPayload.ID.String(),
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    timestamppb.New(refreshPayload.ExpiredAt),
	}
	session, err := createSession(ctx, pool, "https://streamfair_session_service:9093/streamfair/v1/create_session", args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %v.", err)
	}

	rps := &pb_login.LoginUserResponse{
		User:                  ConvertUser(user),
		SessionId:             session.Uuid,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rps, nil
}
