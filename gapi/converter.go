package gapi

import (
	pb_log "github.com/Streamfair/common_proto/IdentityProvider/pb/login"
	"github.com/Streamfair/common_proto/IdentityProvider/pb/register"
	pb "github.com/Streamfair/common_proto/IdentityProvider/pb/register"
	"github.com/Streamfair/common_proto/SessionService/pb/session"
	db "github.com/Streamfair/streamfair_idp/db/sqlc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertSession(session *session.Session) *pb_log.SessionData {
	return &pb_log.SessionData{
		SessionId:        session.GetUuid(),
		UserAgent:        session.GetUserAgent(),
		ClientIp:         session.GetClientIp(),
		IsBlocked:        session.GetIsBlocked(),
		RefreshToken:     session.GetRefreshToken(),
		SessionCreatedAt: session.GetCreatedAt(),
		SessionExpiresAt: session.GetExpiresAt(),
	}
}

func ConvertUserAccount(userAccount db.IdpSvcUserAccount) *register.UserAccount {
	accountData := &pb.AccountData{
		Id:          userAccount.ID,
		AccountName: userAccount.AccountName,
		AccountType: userAccount.AccountType,
		Owner:       userAccount.Owner,
		Bio:         userAccount.Bio,
		Status:      userAccount.Status.String,
		Plan:        userAccount.Plan,
		AvatarUri:   userAccount.AvatarUri.String,
		Plays:       userAccount.Plays,
		Likes:       userAccount.Likes,
		Follows:     userAccount.Follows,
		Shares:      userAccount.Shares,
		CreatedAt:   timestamppb.New(userAccount.AccountCreatedAt),
		UpdatedAt:   timestamppb.New(userAccount.AccountUpdatedAt),
	}

	accountTypeData := &pb.AccountTypeData{
		Id:          userAccount.ID,
		Type:        userAccount.Type,
		Permissions: userAccount.Permissions,
		IsArtist:    userAccount.IsArtist,
		IsProducer:  userAccount.IsProducer,
		IsWriter:    userAccount.IsWriter,
		IsLabel:     userAccount.IsLabel,
		IsUser:      userAccount.IsUser,
		CreatedAt:   timestamppb.New(userAccount.AccountTypeCreatedAt),
		UpdatedAt:   timestamppb.New(userAccount.AccountTypeUpdatedAt),
	}

	userData := &pb.UserData{
		Id:                userAccount.ID,
		Username:          userAccount.Username,
		FullName:          userAccount.FullName,
		Email:             userAccount.Email,
		PasswordHash:      userAccount.PasswordHash,
		PasswordSalt:      userAccount.PasswordSalt,
		CountryCode:       userAccount.CountryCode,
		RoleId:            userAccount.RoleID.Int64,
		Status:            userAccount.Status.String,
		LastLoginAt:       timestamppb.New(userAccount.LastLoginAt),
		UsernameChangedAt: timestamppb.New(userAccount.UsernameChangedAt),
		EmailChangedAt:    timestamppb.New(userAccount.EmailChangedAt),
		PasswordChangedAt: timestamppb.New(userAccount.PasswordChangedAt),
		CreatedAt:         timestamppb.New(userAccount.UserCreatedAt),
		UpdatedAt:         timestamppb.New(userAccount.UserUpdatedAt),
	}

	return &register.UserAccount{
		User:        userData,
		Account:     accountData,
		AccountType: accountTypeData,
	}
}
