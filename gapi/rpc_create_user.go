package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/tpmdigital/simplebank/db/sqlc"
	"github.com/tpmdigital/simplebank/pb"
	"github.com/tpmdigital/simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	// Hash the clear text password
	hashedPassword, err := util.HashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	// Call into the db to create the user
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			// can't create a user with the same email address or username
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}

	// convert dbUser to pbUser
	response := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	// no errors return to the client
	return response, nil
}
