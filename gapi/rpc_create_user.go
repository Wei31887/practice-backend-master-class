package gapi

import (
	"context"
	"fmt"
	db "lesson/simple-bank/db/sqlc"
	pb "lesson/simple-bank/pb"
	"lesson/simple-bank/utils"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPd, err := utils.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Can't hash password: %s", err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPd,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			fmt.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "User already exsit: %s", err.Error())
			}
		}
		return nil, status.Errorf(codes.Internal, "Can't create user: %s", err.Error())
	}

	rsp := &pb.CreateUserResponse{
		User: converUser(user),
	}
	
	return rsp, nil 
}