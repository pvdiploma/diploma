package auth

import (
	"context"

	ssov1 "github.com/pvdiploma/diploma-protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{
		Token: req.GetLogin(),
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	return &ssov1.RegisterResponse{
		UserId: 1,
	}, nil

}

func (s *serverAPI) IsOrginiser(ctx context.Context, req *ssov1.IsOrginiserRequest) (*ssov1.IsOrginiserResponse, error) {
	return &ssov1.IsOrginiserResponse{
		IsOrginiser: true,
	}, nil

}

func (s *serverAPI) IsDistributor(ctx context.Context, req *ssov1.IsDistributorRequest) (*ssov1.IsDistributorResponse, error) {
	return &ssov1.IsDistributorResponse{
		IsDistributor: true,
	}, nil

}

func (s *serverAPI) IsBuyer(ctx context.Context, req *ssov1.IsBuyerRequest) (*ssov1.IsBuyerResponse, error) {
	return &ssov1.IsBuyerResponse{
		IsBuyer: true,
	}, nil

}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	return &ssov1.IsAdminResponse{
		IsAdmin: true,
	}, nil
}
