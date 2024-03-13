package authgrpc

import (
	"context"
	"errors"
	"tn/internal/services/auth"
	"tn/internal/storage"

	ssov1 "github.com/pvdiploma/diploma-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int32) (tokenID string, err error)
	Register(ctx context.Context, login string, email string, password string, role int32) (userID int64, err error)
	IsOrginiser(ctx context.Context, userID int64) (bool, error)
	IsDistributor(ctx context.Context, userID int64) (bool, error)
	IsBuyer(ctx context.Context, userID int64) (bool, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	//TODO: validate data

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	userId, err := s.auth.Register(ctx, req.GetLogin(), req.GetEmail(), req.GetPassword(), req.GetRole())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user with this email arealady exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userId,
	}, nil

}

func (s *serverAPI) IsOrginiser(ctx context.Context, req *ssov1.IsOrginiserRequest) (*ssov1.IsOrginiserResponse, error) {

	isOrginiser, err := s.auth.IsOrginiser(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsOrginiserResponse{
		IsOrginiser: isOrginiser,
	}, nil

}

func (s *serverAPI) IsDistributor(ctx context.Context, req *ssov1.IsDistributorRequest) (*ssov1.IsDistributorResponse, error) {

	isDistributor, err := s.auth.IsDistributor(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsDistributorResponse{
		IsDistributor: isDistributor,
	}, nil

}

func (s *serverAPI) IsBuyer(ctx context.Context, req *ssov1.IsBuyerRequest) (*ssov1.IsBuyerResponse, error) {

	isBuyer, err := s.auth.IsBuyer(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsBuyerResponse{
		IsBuyer: isBuyer,
	}, nil

}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
