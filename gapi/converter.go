package gapi

import (
	"github.com/Streamfair/streamfair_idp/pb/login"
	pb "github.com/Streamfair/streamfair_user_svc/pb/user"
)

func ConvertUser(user *pb.User) *login.Idp_User {
	return &login.Idp_User{
		Id:                user.Id,
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
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
