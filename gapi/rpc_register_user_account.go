package gapi

import (
	"context"
	"encoding/base64"
	"time"

	accountpb "github.com/Streamfair/common_proto/AccountService/pb"
	acc "github.com/Streamfair/common_proto/AccountService/pb/account"
	"github.com/Streamfair/common_proto/AccountService/pb/account_type"
	pb "github.com/Streamfair/common_proto/IdentityProvider/pb/register"
	userpb "github.com/Streamfair/common_proto/UserService/pb"
	"github.com/Streamfair/common_proto/UserService/pb/user"
	"github.com/Streamfair/streamfair_idp/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	IdpSvcAddress       = "streamfair_idp:9091"
	UserSvcAddress      = "streamfair_user_service:9094"
	AccSvcAddress       = "streamfair_account_service:9095"
	TOKEN_svc_address   = "streamfair_token_service:9092"
	SESSION_svc_address = "streamfair_session_service:9093"
)

func (server *Server) RegisterUserAccount(
	ctx context.Context,
	req *pb.RegisterUserAccountRequest) (*pb.RegisterUserAccountResponse, error) {

	pool := NewClientPool(&PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               IdpSvcAddress,
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if userExists(ctx, pool, req.GetUsername()) {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	registeredUser, err := registerNewUser(ctx, pool, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}

	account, accountType, err := createNewAccount(ctx, pool, registeredUser.User.GetUsername(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create account: %v", err)
	}

	return buildRegisterUserAccountResponse(registeredUser, account, accountType), nil
}

func userExists(ctx context.Context, pool *ConnectionPool, username string) bool {
	_, err := getUser(ctx, pool, UserSvcAddress, username)
	return err == nil
}

func registerNewUser(ctx context.Context, pool *ConnectionPool, req *pb.RegisterUserAccountRequest) (*user.CreateUserResponse, error) {
	conn, err := pool.GetConn(UserSvcAddress)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}
	defer conn.Close()

	client := userpb.NewUserServiceClient(conn)

	byteHash, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to prepare registering %v", err)
	}
	hashedPassword := base64.StdEncoding.EncodeToString(byteHash.Hash)
	passwordSalt := base64.StdEncoding.EncodeToString(byteHash.Salt)

	createUserRequest := &user.CreateUserRequest{
		Username:     req.GetUsername(),
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		PasswordSalt: passwordSalt,
		CountryCode:  req.GetCountryCode(),
		RoleId:       int64(req.GetRoleId()),
		Status:       req.GetStatus(),
	}

	return client.CreateUser(ctx, createUserRequest)
}

func createNewAccount(ctx context.Context, pool *ConnectionPool, username string, req *pb.RegisterUserAccountRequest) (*acc.CreateAccountResponse, *account_type.CreateAccountTypeResponse, error) {
	conn, err := pool.GetConn(AccSvcAddress)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to connect to Account Service: %v", err)
	}
	defer conn.Close()

	client := accountpb.NewAccountServiceClient(conn)

	accountTypeExists, err := userOwnsAccountType(ctx, client, username, req.GetAccountType())
	if err != nil {
		return nil, nil, err
	}
	if accountTypeExists {
		return nil, nil, status.Errorf(codes.AlreadyExists, "user already owns account type %d", req.GetAccountType())
	}

	accountType, err := client.CreateAccountType(ctx, buildCreateAccountTypeRequest(req))
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create account type: %v", err)
	}

	account, err := client.CreateAccount(ctx, buildCreateAccountRequest(req, username, accountType.AccountType.Type))
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
	}

	return account, accountType, nil
}

func userOwnsAccountType(ctx context.Context, client accountpb.AccountServiceClient, username string, accountType int32) (bool, error) {
	argList := &acc.ListAccountByOwnerRequest{
		Owner:  username,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := client.ListAccountByOwner(ctx, argList)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, nil
		}
		return false, status.Errorf(codes.Internal, "failed to check accounts for username: %v", err)
	}

	for _, account := range accounts.GetAccounts() {
		if account.GetAccountType() == accountType {
			return true, nil
		}
	}
	return false, nil
}

func buildCreateAccountTypeRequest(req *pb.RegisterUserAccountRequest) *account_type.CreateAccountTypeRequest {
	return &account_type.CreateAccountTypeRequest{
		Type:        req.GetAccountType(),
		Permissions: "default_user",
		IsArtist:    false,
		IsProducer:  false,
		IsWriter:    false,
		IsLabel:     false,
		IsUser:      true,
	}
}

func buildCreateAccountRequest(req *pb.RegisterUserAccountRequest, username string, accountType int32) *acc.CreateAccountRequest {
	return &acc.CreateAccountRequest{
		AccountType: accountType,
		AccountName: req.GetAccountName(),
		Owner:       username,
		Bio:         req.GetAccountBio(),
		Status:      "active",
		Plan:        req.GetAccountPlan(),
		AvatarUri:   req.GetAccountAvatarUri(),
		Plays:       0,
		Likes:       0,
		Follows:     0,
		Shares:      0,
	}
}

func buildRegisterUserAccountResponse(userResp *user.CreateUserResponse, accountResp *acc.CreateAccountResponse, accountTypeResp *account_type.CreateAccountTypeResponse) *pb.RegisterUserAccountResponse {
accountData := &pb.AccountData{
	Id:          accountResp.Account.Id,
	AccType:     accountResp.Account.AccountType,
	AccountName: accountResp.Account.AccountName,
	Owner:       accountResp.Account.Owner,
	Bio:         accountResp.Account.Bio,
	Status:      accountResp.Account.Status,
	Plan:        accountResp.Account.Plan,
	AvatarUri:   accountResp.Account.AvatarUri,
	Plays:       accountResp.Account.Plays,
	Likes:       accountResp.Account.Likes,
	Follows:     accountResp.Account.Follows,
	Shares:      accountResp.Account.Shares,
	CreatedAt:   accountResp.Account.CreatedAt,
	UpdatedAt:   accountResp.Account.UpdatedAt,
	AccountType: &pb.AccountTypeData{
		Id:          accountTypeResp.AccountType.Id,
		AccType:     accountTypeResp.AccountType.Type,
		Permissions: accountTypeResp.AccountType.Permissions,
		IsArtist:    accountTypeResp.AccountType.IsArtist,
		IsProducer:  accountTypeResp.AccountType.IsProducer,
		IsWriter:    accountTypeResp.AccountType.IsWriter,
		IsLabel:     accountTypeResp.AccountType.IsLabel,
		IsUser:      accountTypeResp.AccountType.IsUser,
		CreatedAt:   accountTypeResp.AccountType.CreatedAt,
		UpdatedAt:   accountTypeResp.AccountType.UpdatedAt,
	},
}

userData := &pb.UserData{
	Id:                userResp.User.Id,
	Username:          userResp.User.Username,
	FullName:          userResp.User.FullName,
	Email:             userResp.User.Email,
	PasswordHash:      userResp.User.PasswordHash,
	PasswordSalt:      userResp.User.PasswordSalt,
	CountryCode:       userResp.User.CountryCode,
	RoleId:            userResp.User.RoleId,
	Status:            userResp.User.Status,
	AccountType:       accountData.AccType,
	AccountName:       accountData.AccountName,
	AccountBio:        accountData.Bio,
	AccountPlan:       accountData.Plan,
	AccountAvatarUri:  accountData.AvatarUri,
	LastLoginAt:       userResp.User.LastLoginAt,
	UsernameChangedAt: userResp.User.UsernameChangedAt,
	EmailChangedAt:    userResp.User.EmailChangedAt,
	PasswordChangedAt: userResp.User.PasswordChangedAt,
	CreatedAt:         userResp.User.CreatedAt,
	UpdatedAt:         userResp.User.UpdatedAt,
}

	return &pb.RegisterUserAccountResponse{
		UserAccount: &pb.UserAccount{
			User:    userData,
			Account: accountData,
		},
	}
}