package main

import (
	"context"
	pb "ordermgt/service/ecommerce"

	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	orderMap map[string]*pb.Order
}

func (s *server) GetOrder(ctx context.Context,
	orderId *wrapper.StringValue) (*pb.Order, error) {
	ord, exists := s.orderMap[orderId.Value]
	if exists {
		return ord, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "Order does not exist. : ", orderId)
}
