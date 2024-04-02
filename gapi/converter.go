package gapi

// func ConvertUser(user *pb.User) *login.Idp_User {
// 	return &login.Idp_User{
// 		Id:                user.Id,
// 		Username:          user.Username,
// 		FullName:          user.FullName,
// 		Email:             user.Email,
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

// func ConvertUser(user *user.User) *pb.User {
// 	return &pb.User{
// 		Id:                user.ID,
// 		Username:          user.Username,
// 		FullName:          user.FullName,
// 		Email:             user.Email,
// 		CountryCode:       user.CountryCode,
// 		RoleId:            user.RoleID.Int64,
// 		Status:            user.Status.String,
// 		LastLoginAt:       timestamppb.New(user.LastLoginAt),
// 		UsernameChangedAt: timestamppb.New(user.UsernameChangedAt),
// 		EmailChangedAt:    timestamppb.New(user.EmailChangedAt),
// 		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
// 		CreatedAt:         timestamppb.New(user.CreatedAt),
// 		UpdatedAt:         timestamppb.New(user.UpdatedAt),
// 	}
// }
