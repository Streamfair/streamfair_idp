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
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	IdpSvcAddress  = "streamfair_idp:9091"
	UserSvcAddress = "streamfair_user_service:9094"
	AccSvcAddress  = "streamfair_account_service:9095"
)

var (
	registered         bool
	accountTypeCreated bool
	accountCreated     bool
)

type ExtendedRegisterUserAccount struct {
	*pb.RegisterUserAccountResponse `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Account                         *structpb.Struct `protobuf:"bytes,100,opt,name=account,proto3" json:"account,omitempty"`
}

// RegisterUserAccount involves: UserService
// 1. Creates the user with the given information via UserService.
// 2. If the user already exists, return an error.
// 3. Hash the password and store it in the database.
// 4. Create user account with the given information via AccountService.
// 5. Return the user account information.
func (server *Server) RegisterUserAccount(
	ctx context.Context,
	req *pb.RegisterUserAccountRequest) (*pb.RegisterUserAccountResponse, error) {
	poolConfig := &PoolConfig{
		MaxOpenConnection:     10,
		MaxIdleConnection:     5,
		ConnectionQueueLength: 10,
		Address:               IdpSvcAddress,
		ConfigOptions:         []grpc.DialOption{},
		IdleTimeout:           10 * time.Second,
	}

	pool := NewClientPool(poolConfig)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	username := req.GetUsername()
	_, err := getUser(ctx, pool, UserSvcAddress, username)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	registered_user, err := registerUser(ctx, pool, UserSvcAddress, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
	}
	registered = true

	reqacctype := &account_type.CreateAccountTypeRequest{
		Type:        req.GetAccountType(),
		Permissions: "default",
		IsArtist:    false,
		IsProducer:  false,
		IsWriter:    false,
		IsLabel:     false,
		IsUser:      true,
	}

	reqacc := &acc.CreateAccountRequest{
		AccountType: reqacctype.GetType(),
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

	account, accountType, err := createAccount(
		ctx,
		pool,
		AccSvcAddress,
		reqacc,
		reqacctype,
		registered_user.User.GetUsername())
	if err != nil {
		// if (registered) {
		// 	// Rollback user registration
		// 	_, err := deleteUser(ctx, pool, USER_svc_address, username)
		// 	if err != nil {
		// 		return nil, status.Errorf(codes.Internal, "Failed to rollback user registration: %v", err)
		// 	}
		// }
		return nil, status.Errorf(codes.Internal, "Failed to create account: %v", err)
	}

	userAccountData, err := structpb.NewStruct(map[string]interface{}{
		"accountType": accountType,
		"account":     account,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create account: %v", err)
	}

	rsp := &pb.RegisterUserAccountResponse{
		UserAccount: &pb.UserAccount{
			Username:          registered_user.User.GetUsername(),
			FullName:          registered_user.User.GetFullName(),
			Email:             registered_user.User.GetEmail(),
			PasswordHash:      registered_user.User.GetPasswordHash(),
			PasswordSalt:      registered_user.User.GetPasswordSalt(),
			CountryCode:       registered_user.User.GetCountryCode(),
			RoleId:            registered_user.User.GetRoleId(),
			Status:            registered_user.User.GetStatus(),
			AccountType:       req.GetAccountType(),
			AccountName:       req.GetAccountName(),
			AccountBio:        req.GetAccountBio(),
			AccountPlan:       req.GetAccountPlan(),
			AccountAvatarUri:  req.GetAccountAvatarUri(),
			LastLoginAt:       registered_user.User.GetLastLoginAt(),
			UsernameChangedAt: registered_user.User.GetUsernameChangedAt(),
			EmailChangedAt:    registered_user.User.GetEmailChangedAt(),
			PasswordChangedAt: registered_user.User.GetPasswordChangedAt(),
			CreatedAt:         registered_user.User.GetCreatedAt(),
			UpdatedAt:         registered_user.User.GetUpdatedAt(),
		},
	}

	extendedRsp := &ExtendedRegisterUserAccount{
		RegisterUserAccountResponse: rsp,
		Account:                     userAccountData,
	}

	// Convert RegisterUserAccount back to *pb.RegisterUserResponse for the return value
	userAccount := &pb.RegisterUserAccountResponse{
		UserAccount: extendedRsp.UserAccount,
	}
	return userAccount, nil
}

func createAccount(
	ctx context.Context,
	pool *ConnectionPool,
	address string,
	reqacc *acc.CreateAccountRequest,
	reqacctype *account_type.CreateAccountTypeRequest,
	username string) (*acc.CreateAccountResponse, *account_type.CreateAccountTypeResponse, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to connect to Account Service: %v", err)
	}

	client := accountpb.NewAccountServiceClient(conn)

	accounts, err := client.ListAccountByOwner(ctx, &acc.ListAccountByOwnerRequest{
		Owner:  username,
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to list accounts: %v", err)
	}

	if len(accounts.Accounts) > 0 {
		return nil, nil, nil
	}

	accountType, err := client.CreateAccountType(ctx, &account_type.CreateAccountTypeRequest{
		Type:        reqacctype.GetType(),
		Permissions: reqacctype.GetPermissions(),
		IsArtist:    reqacctype.GetIsArtist(),
		IsProducer:  reqacctype.GetIsProducer(),
		IsWriter:    reqacctype.GetIsWriter(),
		IsLabel:     reqacctype.GetIsLabel(),
		IsUser:      reqacctype.GetIsUser(),
	})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create default account type: %v", err)
	}
	accountTypeCreated = true

	arg := &acc.CreateAccountRequest{
		AccountType: accountType.AccountType.GetType(),
		AccountName: reqacc.GetAccountName(),
		Owner:       reqacc.GetOwner(),
		Bio:         reqacc.GetBio(),
		Status:      reqacc.GetStatus(),
		Plan:        reqacc.GetPlan(),
		AvatarUri:   reqacc.GetAvatarUri(),
		Plays:       reqacc.GetPlays(),
		Likes:       reqacc.GetLikes(),
		Follows:     reqacc.GetFollows(),
		Shares:      reqacc.GetShares(),
	}

	account, err := client.CreateAccount(ctx, arg)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create default account: %v", err)
	}
	accountCreated = true

	return account, accountType, nil
}

func registerUser(
	ctx context.Context,
	pool *ConnectionPool,
	address string,
	req *pb.RegisterUserAccountRequest) (*user.CreateUserResponse, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
	}

	client := userpb.NewUserServiceClient(conn)

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

	return resp, nil
}

// func deleteUser(ctx context.Context, pool *ConnectionPool, address string, req *pb.RegisterUserAccountRequest) (*user.CreateUserResponse, error) {
// 	conn, err := pool.GetConn(address)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to connect to UserService: %v", err)
// 	}

// 	client := userpb.NewUserServiceClient(conn)

// 	byteHash, err := util.HashPassword(req.Password)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to prepare registering %v", err)
// 	}
// 	hashedPassword := base64.StdEncoding.EncodeToString(byteHash.Hash)
// 	passwordSalt := base64.StdEncoding.EncodeToString(byteHash.Salt)

// 	arg := &user.CreateUserRequest{
// 		Username:     req.GetUsername(),
// 		FullName:     req.GetFullName(),
// 		Email:        req.GetEmail(),
// 		PasswordHash: hashedPassword,
// 		PasswordSalt: passwordSalt,
// 		CountryCode:  req.GetCountryCode(),
// 		RoleId:       int64(req.GetRoleId()),
// 		Status:       req.GetStatus(),
// 	}
// 	resp, err := client.CreateUser(ctx, arg)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
// 	}

// 	return resp, nil
// }
