package gapi

import (
	"context"
	"encoding/base64"
	"time"

	account_pb "github.com/Streamfair/common_proto/AccountService/pb"
	account "github.com/Streamfair/common_proto/AccountService/pb/account"
	account_type "github.com/Streamfair/common_proto/AccountService/pb/account_type"
	pb_register "github.com/Streamfair/common_proto/IdentityProvider/pb/register"
	user_pb "github.com/Streamfair/common_proto/UserService/pb"
	user "github.com/Streamfair/common_proto/UserService/pb/user"
	"github.com/Streamfair/streamfair_idp/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	IDP_svc_address  = "streamfair_idp:9091"
	USER_svc_address = "streamfair_user_service:9094"
	ACC_svc_address  = "streamfair_account_service:9095"
)

// RegisterUserAccount registers a user in the usersdb, check if it has an account and if not, create a default account.
// Returns both user and account.

// Microservices involved: UserService
// 1. Creates the user with the given information via UserService.
// 2. If the user already exists, return an error.
// 3. Hash the password and store it in the database.
// 4. If not existing: create user account with the given information via AccountService.
// 5. Return the user and account information.
type RegisterUserAccount struct {
    *pb_register.RegisterUserResponse `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
    AdditionalData *structpb.Struct `protobuf:"bytes,100,opt,name=additional_data,proto3" json:"additional_data,omitempty"`
}

func (server *Server) RegisterUserAccount(ctx context.Context, req *pb_register.RegisterUserRequest) (*pb_register.RegisterUserResponse, error) {
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

    _, err := getUser(ctx, pool, USER_svc_address, username)
    if err == nil {
        return nil, status.Errorf(codes.AlreadyExists, "user already exists")
    }

    registered_user, err := registerUser(ctx, pool, USER_svc_address, req)
    if err!= nil {
        return nil, status.Errorf(codes.Internal, "Failed to register user: %v", err)
    }

	req_acc_type := &account_type.CreateAccountTypeRequest{
		Type:        1,
		Permissions: "default",
		IsArtist:    false,
		IsProducer:  false,
		IsWriter:    false,
		IsLabel:     false,
		IsUser:      true,
	}
	
	req_acc := &account.CreateAccountRequest{
		AccountType: req_acc_type.GetType(),
		AccountName: username,
		Owner:       username,
		Bio:         "This is a default account",
		Status:      "active",
		Plan:        "free",
		AvatarUri:   "default",
		Plays:       0,
		Likes:       0,
		Follows:     0,
		Shares:      0,
	}

    account, accountType, err := createAccount(ctx, pool, ACC_svc_address, req_acc, req_acc_type, registered_user.User.GetUsername())
    if err!= nil {
        return nil, status.Errorf(codes.Internal, "Failed to create account: %v", err)
    }

    userAccountStruct, _ := structpb.NewStruct(map[string]interface{}{
        "accountType": accountType,
        "account":     account,
    })

    rsp := &pb_register.RegisterUserResponse{
        User: &pb_register.User{
            Username:          registered_user.User.GetUsername(),
            FullName:          registered_user.User.GetFullName(),
            Email:             registered_user.User.GetEmail(),
            PasswordHash:      registered_user.User.GetPasswordHash(),
            PasswordSalt:      registered_user.User.GetPasswordSalt(),
            CountryCode:       registered_user.User.GetCountryCode(),
            RoleId:            registered_user.User.GetRoleId(),
            Status:            registered_user.User.GetStatus(),
            LastLoginAt:       registered_user.User.GetLastLoginAt(),
            UsernameChangedAt: registered_user.User.GetUsernameChangedAt(),
            EmailChangedAt:    registered_user.User.GetEmailChangedAt(),
            PasswordChangedAt: registered_user.User.GetPasswordChangedAt(),
            CreatedAt:         registered_user.User.GetCreatedAt(),
            UpdatedAt:         registered_user.User.GetUpdatedAt(),
        },
    }

    extendedRsp := &RegisterUserAccount{
        RegisterUserResponse: rsp,
        AdditionalData:        userAccountStruct,
    }

    // Convert RegisterUserAccount back to *pb_register.RegisterUserResponse for the return value
    regResp := &pb_register.RegisterUserResponse{
        User: extendedRsp.User,
    }
    return regResp, nil
}

func createAccount(ctx context.Context, pool *ConnectionPool, address string, req_acc *account.CreateAccountRequest, req_acc_type *account_type.CreateAccountTypeRequest, username string) (*account.CreateAccountResponse, *account_type.CreateAccountTypeResponse, error) {
	conn, err := pool.GetConn(address)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to connect to Account Service: %v", err)
	}

	client := account_pb.NewAccountServiceClient(conn)

	accounts, err := client.ListAccountByOwner(ctx, &account.ListAccountByOwnerRequest{
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
		Type:        req_acc_type.GetType(),
		Permissions: req_acc_type.GetPermissions(),
		IsArtist:    req_acc_type.GetIsArtist(),
		IsProducer:  req_acc_type.GetIsProducer(),
		IsWriter:    req_acc_type.GetIsWriter(),
		IsLabel:     req_acc_type.GetIsLabel(),
		IsUser:      req_acc_type.GetIsUser(),
	})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create default account type: %v", err)
	}

	arg := &account.CreateAccountRequest{
		AccountType: accountType.AccountType.GetType(),
		AccountName: req_acc.GetAccountName(),
		Owner:       req_acc.GetOwner(),
		Bio:         req_acc.GetBio(),
		Status:      req_acc.GetStatus(),
		Plan:        req_acc.GetPlan(),
		AvatarUri:   req_acc.GetAvatarUri(),
		Plays:       req_acc.GetPlays(),
		Likes:       req_acc.GetLikes(),
		Follows:     req_acc.GetFollows(),
		Shares:      req_acc.GetShares(),
	}

	account, err := client.CreateAccount(ctx, arg)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create default account: %v", err)
	}

	return account, accountType, nil
}

func registerUser(ctx context.Context, pool *ConnectionPool, address string, req *pb_register.RegisterUserRequest) (*user.CreateUserResponse, error) {
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

	return resp, nil
}
