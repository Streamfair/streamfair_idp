package gapi

import (
	pb_log "github.com/Streamfair/common_proto/IdentityProvider/pb/login"
	user "github.com/Streamfair/common_proto/UserService/pb/user"
)

// func ConvertRegisteredUser(user *user.User) *pb_reg.UserAccount {
// 	return &pb_reg.UserAccount{
// 		Id:                user.Id,
// 		Username:          user.Username,
// 		FullName:          user.FullName,
// 		Email:             user.Email,
// 		PasswordHash:      user.PasswordHash,
// 		PasswordSalt:      user.PasswordSalt,
// 		CountryCode:       user.CountryCode,
// 		RoleId:            user.RoleId,
// 		Status:            user.Status,
// 		LastLoginAt:       user.LastLoginAt,
// 		UsernameChangedAt: user.UsernameChangedAt,
// 		EmailChangedAt:    user.EmailChangedAt,
// 		PasswordChangedAt: user.PasswordChangedAt,
// 		CreatedAt:         user.CreatedAt,
// 		UpdatedAt:         user.UpdatedAt,
// 	}
// }

func ConvertLoggedInUserAcount(user *user.User) *pb_log.LoggedInUser {
	return &pb_log.LoggedInUser{
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
