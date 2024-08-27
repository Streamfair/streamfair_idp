package gapi

import (
	"context"

	"google.golang.org/grpc/codes"
	"github.com/Streamfair/streamfair_idp/validator"
	pb "github.com/Streamfair/common_proto/IdentityProvider/pb/register"

)

func (server *Server) GetUserAccountByUsername(ctx context.Context, req *pb.GetUserAccountByUsernameRequest) (*pb.GetUserAccountByUsernameResponse, error) {
	username := req.GetUsername()

	// Perform field validation
	err := validator.ValidateUsername(username)
	if err != nil {
		violation := (&CustomError{
			StatusCode: codes.InvalidArgument,
		}).WithDetails("username", err)
		return nil, invalidArgumentError(violation)
	}

	// Fetch user from the database
	userAccount, err := server.store.GetUserAccountByUserAccountname(ctx, username)
	if err != nil {
		return nil, handleDatabaseError(err)
	}

	rsp := &pb.GetUserAccountByUsernameResponse{
		UserAccount: ConvertUserAccount(userAccount),
	}
	return rsp, nil
}
