package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	user "github.com/anhhuy1010/cms-user/grpc/proto/users"

	"github.com/anhhuy1010/cms-user/models"
)

type UserService struct {
}

func NewUserServer() user.UserServer {
	return &UserService{}
}

// Detail implements user.UserServer.
func (s *UserService) Detail(ctx context.Context, req *user.DetailRequest) (*user.DetailResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "Token is required")
	}

	conditions := bson.M{"token": req.Token}

	result, err := new(models.Tokens).FindOne(conditions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "Token not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	cond := bson.M{"uuid": result.UserUuid}
	userResult, err := new(models.Users).FindOne(cond)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	res := &user.DetailResponse{
		UserUuid: result.UserUuid,
		Role:     userResult.Role,
	}
	return res, nil
}

// func (s *UserService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
// 	conditions := bson.M{}

// 	result, err := new(models.Users).Find(conditions)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}
// 	var respData []*pb.DetailResponse
// 	for _, user := range result {

// 		res := &pb.DetailResponse{
// 			Uuid:     user.Uuid,
// 			Username: user.Username,
// 			IsActive: int32(user.IsActive),
// 		}
// 		respData = append(respData, res)
// 	}

// 	return &pb.ListResponse{
// 		Users: respData,
// 	}, nil
// }
