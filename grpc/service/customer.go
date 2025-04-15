package service

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	customer "github.com/anhhuy1010/DATN-cms-customer/grpc/proto/customer"
	"github.com/anhhuy1010/DATN-cms-customer/models"
)

type CustomerService struct {
}

func NewCustomerServer() customer.CustomerServiceServer {
	return &CustomerService{}
}

func (s *CustomerService) ListCustomers(ctx context.Context, req *customer.ListCustomerRequest) (*customer.ListCustomerResponse, error) {
	customerModel := models.Customer{}

	cond := bson.M{}
	if req.Username != nil {
		cond["username"] = req.GetUsername()
	}
	if req.IsActive != nil {
		cond["is_active"] = req.GetIsActive()
	}

	optionsQuery, _, _ := models.GetPagingOption(int(req.Page), int(req.Limit), "")
	users, err := customerModel.Pagination(ctx, cond, optionsQuery)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get customers")
	}

	var customers []*customer.Customer
	for _, u := range users {
		customers = append(customers, &customer.Customer{
			Uuid:     u.Uuid,
			Username: u.Username,
			IsActive: int32(u.IsActive),
		})
	}

	return &customer.ListCustomerResponse{
		Data: customers,
	}, nil
}
