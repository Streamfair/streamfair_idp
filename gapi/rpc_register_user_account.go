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
	db "github.com/Streamfair/streamfair_idp/db/sqlc"
	"github.com/Streamfair/streamfair_idp/util"
	"github.com/jackc/pgx/v5/pgtype"
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

	response := buildRegisterUserAccountResponse(registeredUser, account, accountType)

	err = storeUserAccount(ctx, server, response)

	return response, nil
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
		AccountName: req.GetAccountName(),
		AccountType: accountType,
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
		Id:          accountResp.Account.GetId(),
		AccountName: accountResp.Account.GetAccountName(),
		AccountType: accountResp.Account.GetAccountType(),
		Owner:       accountResp.Account.GetOwner(),
		Bio:         accountResp.Account.GetBio(),
		Status:      accountResp.Account.GetStatus(),
		Plan:        accountResp.Account.GetPlan(),
		AvatarUri:   accountResp.Account.GetAvatarUri(),
		Plays:       accountResp.Account.GetPlays(),
		Likes:       accountResp.Account.GetLikes(),
		Follows:     accountResp.Account.GetFollows(),
		Shares:      accountResp.Account.GetShares(),
		CreatedAt:   accountResp.Account.GetCreatedAt(),
		UpdatedAt:   accountResp.Account.GetUpdatedAt(),
	}

	accountTypeData := &pb.AccountTypeData{
		Id:          accountTypeResp.AccountType.GetId(),
		Type:        accountTypeResp.AccountType.GetType(),
		Permissions: accountTypeResp.AccountType.GetPermissions(),
		IsArtist:    accountTypeResp.AccountType.GetIsArtist(),
		IsProducer:  accountTypeResp.AccountType.GetIsProducer(),
		IsWriter:    accountTypeResp.AccountType.GetIsWriter(),
		IsLabel:     accountTypeResp.AccountType.GetIsLabel(),
		IsUser:      accountTypeResp.AccountType.GetIsUser(),
		CreatedAt:   accountTypeResp.AccountType.GetCreatedAt(),
		UpdatedAt:   accountTypeResp.AccountType.GetUpdatedAt(),
	}

	userData := &pb.UserData{
		Id:                userResp.User.GetId(),
		Username:          userResp.User.GetUsername(),
		FullName:          userResp.User.GetFullName(),
		Email:             userResp.User.GetEmail(),
		PasswordHash:      userResp.User.GetPasswordHash(),
		PasswordSalt:      userResp.User.GetPasswordSalt(),
		CountryCode:       userResp.User.GetCountryCode(),
		RoleId:            userResp.User.GetRoleId(),
		Status:            userResp.User.GetStatus(),
		LastLoginAt:       userResp.User.GetLastLoginAt(),
		UsernameChangedAt: userResp.User.GetUsernameChangedAt(),
		EmailChangedAt:    userResp.User.GetEmailChangedAt(),
		PasswordChangedAt: userResp.User.GetPasswordChangedAt(),
		CreatedAt:         userResp.User.GetCreatedAt(),
		UpdatedAt:         userResp.User.GetUpdatedAt(),
	}

	userAccount := &pb.UserAccount{
		User:        userData,
		Account:     accountData,
		AccountType: accountTypeData,
	}

	return &pb.RegisterUserAccountResponse{
		UserAccount: userAccount,
	}
}
func storeUserAccount(ctx context.Context, server *Server, response *pb.RegisterUserAccountResponse) error {
	user := response.GetUserAccount().GetUser()
	account := response.GetUserAccount().GetAccount()
	accountType := response.GetUserAccount().GetAccountType()
	userAccountEntry := db.CreateUserAccountParams{
		Username:             user.GetUsername(),
		FullName:             user.GetFullName(),
		Email:                user.GetEmail(),
		PasswordHash:         user.GetPasswordHash(),
		PasswordSalt:         user.GetPasswordSalt(),
		CountryCode:          user.GetCountryCode(),
		RoleID:               pgtype.Int8{Int64: user.GetRoleId()},
		Status:               pgtype.Text{String: user.GetStatus()},
		LastLoginAt:          user.GetLastLoginAt().AsTime(),
		UsernameChangedAt:    user.GetUsernameChangedAt().AsTime(),
		EmailChangedAt:       user.GetEmailChangedAt().AsTime(),
		PasswordChangedAt:    user.GetPasswordChangedAt().AsTime(),
		UserCreatedAt:        user.GetCreatedAt().AsTime(),
		UserUpdatedAt:        user.GetUpdatedAt().AsTime(),
		AccountName:          account.GetAccountName(),
		AccountType:          account.GetAccountType(),
		Owner:                account.GetOwner(),
		Bio:                  account.GetBio(),
		AccountStatus:        account.GetStatus(),
		Plan:                 account.GetPlan(),
		AvatarUri:            pgtype.Text{String: account.GetAvatarUri()},
		Plays:                account.GetPlays(),
		Likes:                account.GetLikes(),
		Follows:              account.GetFollows(),
		Shares:               account.GetShares(),
		AccountCreatedAt:     account.GetCreatedAt().AsTime(),
		AccountUpdatedAt:     account.GetUpdatedAt().AsTime(),
		Type:                 accountType.GetType(),
		Permissions:          accountType.GetPermissions(),
		IsArtist:             accountType.GetIsArtist(),
		IsProducer:           accountType.GetIsProducer(),
		IsWriter:             accountType.GetIsWriter(),
		IsLabel:              accountType.GetIsLabel(),
		IsUser:               accountType.GetIsUser(),
		AccountTypeCreatedAt: accountType.GetCreatedAt().AsTime(),
		AccountTypeUpdatedAt: accountType.GetUpdatedAt().AsTime(),
	}

	_, err := server.store.CreateUserAccount(ctx, userAccountEntry)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to save user account: %v", err)
	}
	return nil
}

func getUser(ctx context.Context, pool *ConnectionPool, address string, username string) (*user.User, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}

	client := userpb.NewUserServiceClient(conn)

	req := &user.GetUserByValueRequest{
		Username: username,
	}
	resp, err := client.GetUserByValue(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	return resp.User, nil
}