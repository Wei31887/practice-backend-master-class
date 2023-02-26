package gapi

import (
	"context"
	"database/sql"
	db "lesson/simple-bank/db/sqlc"
	pb "lesson/simple-bank/pb"
	"lesson/simple-bank/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Internal error")
	}

	err = utils.ComparePassword(user.HashedPassword, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Password not matched")
	}

	// create the access token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Acces token generation failed")
	}

	// create the refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Refresh token generation failed")
	}

	// store refresh token into session
	arg := db.CreateSessionParams{
		ID:           refreshPayload.Id,
		Username:     req.Username,
		RefreshToken: refreshToken,
		ClientIp:     "",
		UserAgent:    "",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	}

	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Create session failed")
	}

	// send the token to the client
	rsp := &pb.LoginUserResponse{
		SessionId:        session.ID.String(),
		AccessToken:      accessToken,
		AccessExpiredAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshToken:     refreshToken,
		RefreshExpiredAt: timestamppb.New(refreshPayload.ExpiresAt),
		User:             converUser(user),
	}

	return rsp, nil 
}
