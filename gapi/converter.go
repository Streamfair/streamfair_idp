package gapi

import (
	pb_reg "github.com/Streamfair/common_proto/IdentityProvider/pb/register"
	pb_log "github.com/Streamfair/common_proto/IdentityProvider/pb/login"
	user "github.com/Streamfair/common_proto/UserService/pb/user"
)

func ConvertRegisteredUser(user *user.User) *pb_reg.UserRegistered {
	return &pb_reg.UserRegistered{
		Id:                user.Id,
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordHash:      user.PasswordHash,
		PasswordSalt:      user.PasswordSalt,
		CountryCode:       user.CountryCode,
		RoleId:            user.RoleId,
		Status:            user.Status,
		LastLoginAt:       user.LastLoginAt,
		UsernameChangedAt: user.UsernameChangedAt,
		EmailChangedAt:    user.EmailChangedAt,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}

func ConvertLoggedInUser(user *user.User) *pb_log.UserLogin {
	return &pb_log.UserLogin{
		Id:                user.Id,
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordHash:      user.PasswordHash,
		PasswordSalt:      user.PasswordSalt,
		CountryCode:       user.CountryCode,
		RoleId:            user.RoleId,
		Status:            user.Status,
		LastLoginAt:       user.LastLoginAt,
		UsernameChangedAt: user.UsernameChangedAt,
		EmailChangedAt:    user.EmailChangedAt,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}
}
