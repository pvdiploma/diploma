package auth

import (
	"context"
	"log/slog"
	"time"
	"tn/internal/domain/models"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int32) (tokenID string, err error)
	Register(ctx context.Context, login string, email string, password string, role int32) (userID int64, err error)
	IsOrginiser(ctx context.Context, userID int64) (bool, error)
	IsDistributor(ctx context.Context, userID int64) (bool, error)
	IsBuyer(ctx context.Context, userID int64) (bool, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type UserStorage interface {
	SaveUser(ctx context.Context, login string, email string, pwdHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsOrginiser(ctx context.Context, userID int64) (bool, error)
	IsDistributor(ctx context.Context, userID int64) (bool, error)
	IsBuyer(ctx context.Context, userID int64) (bool, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type AuthService struct {
	log          *slog.Logger
	userStorage  UserStorage
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

// New return an instance of auth service
func New(
	log *slog.Logger,
	userStorage UserStorage,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,

) *AuthService {
	return &AuthService{
		log:          log,
		userStorage:  userStorage,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, email string, password string, appID int32) (string, error) {
	panic("not emplemented")
}

func (a *AuthService) Register(ctx context.Context, login string, email string, password string, role uint32, appID int32) (string, error) {
	panic("not implemented")
}

func (a *AuthService) IsOrginiser(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}

func (a *AuthService) IsDistributor(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}

func (a *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}
