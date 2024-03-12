package auth

import (
	"context"
	"tn/internal/services/auth"

	ssov1 "github.com/pvdiploma/diploma-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth auth.Auth
}

func Register(gRPC *grpc.Server, auth auth.Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	//TODO: validate data

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword(), req.GetAppId())
	if err != nil {
		//TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	userId, err := s.auth.Register(ctx, req.GetLogin(), req.GetEmail(), req.GetPassword(), req.GetRole())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil

}

func (s *serverAPI) IsOrginiser(ctx context.Context, req *ssov1.IsOrginiserRequest) (*ssov1.IsOrginiserResponse, error) {

	isOrginiser, err := s.auth.IsOrginiser(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsOrginiserResponse{
		IsOrginiser: isOrginiser,
	}, nil

}

func (s *serverAPI) IsDistributor(ctx context.Context, req *ssov1.IsDistributorRequest) (*ssov1.IsDistributorResponse, error) {

	isDistributor, err := s.auth.IsDistributor(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsDistributorResponse{
		IsDistributor: isDistributor,
	}, nil

}

func (s *serverAPI) IsBuyer(ctx context.Context, req *ssov1.IsBuyerRequest) (*ssov1.IsBuyerResponse, error) {

	isBuyer, err := s.auth.IsBuyer(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsBuyerResponse{
		IsBuyer: isBuyer,
	}, nil

}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
