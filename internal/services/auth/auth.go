package auth

import (
	"context"
	"errors"
	"log/slog"
	"tn/internal/domain/models"
	"tn/internal/storage"
	tokenmanager "tn/internal/utils/tokenManager"
	sl "tn/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	SaveUser(ctx context.Context, login string, email string, pwdHash []byte, role int32) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsOrginiser(ctx context.Context, userID int64) (bool, error)
	IsDistributor(ctx context.Context, userID int64) (bool, error)
	IsBuyer(ctx context.Context, userID int64) (bool, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int32) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid appID")
)

type AuthService struct {
	log          *slog.Logger
	userStorage  UserStorage
	userProvider UserProvider
	appProvider  AppProvider
	tokenManager *tokenmanager.TokenManager
}

// New return an instance of auth service
func New(
	log *slog.Logger,
	userStorage UserStorage,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenManager *tokenmanager.TokenManager,

) *AuthService {
	return &AuthService{
		log:          log,
		userStorage:  userStorage,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenManager: tokenManager,
	}
}

func (a *AuthService) Login(ctx context.Context, email string, password string, appID int32) (string, error) {

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return "", ErrInvalidCredentials
		}

		a.log.Error("failed to get user info", sl.Err(err))
	}

	err = bcrypt.CompareHashAndPassword(user.PwdHash, []byte(password))
	if err != nil {
		a.log.Warn("password compare error", sl.Err(err))
		return "", ErrInvalidCredentials
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		a.log.Warn("app not found", sl.Err(err))
		return "", err
	}
	_ = app

	tokens, err := a.tokenManager.GenerateNewJWTPair(user)

	if err != nil {
		a.log.Error("refresh token generation error", sl.Err(err))
		return "", ErrInvalidCredentials
	}

	a.log.Info("user logged in successfuly")
	// put tokens.RefreshToken to redis or mongodb
	return tokens.AccessToken, nil
}

func (a *AuthService) Register(ctx context.Context, login string, email string, password string, role int32, appID int32) (int64, error) {

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		a.log.Error(" password hash deneration error", sl.Err(err))
		return -1, err
	}
	_ = pwdHash

	id, err := a.userStorage.SaveUser(ctx, login, email, pwdHash, role)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("user already exists", sl.Err(err))
			return -1, ErrInvalidCredentials
		}
		a.log.Error(" saving user error", sl.Err(err))
		return -1, err
	}

	// call login??
	_, err = a.Login(ctx, email, password, appID)
	return id, nil
}

func (a *AuthService) IsOrginiser(ctx context.Context, userID int64) (bool, error) {
	IsOrginiser, err := a.userProvider.IsOrginiser(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", sl.Err(err))
			return false, ErrInvalidAppID
		}
		return false, err
	}

	a.log.Info("checking if user is orginiser", slog.Int64("userID", userID), slog.Bool("isOrginiser", IsOrginiser))
	return IsOrginiser, nil
}

func (a *AuthService) IsDistributor(ctx context.Context, userID int64) (bool, error) {
	IsDistributor, err := a.userProvider.IsDistributor(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", sl.Err(err))
			return false, ErrInvalidAppID
		}
		return false, err
	}

	a.log.Info("checking if user is distributor", slog.Int64("userID", userID), slog.Bool("isDistributor", IsDistributor))
	return IsDistributor, nil
}

func (a *AuthService) IsBuyer(ctx context.Context, userID int64) (bool, error) {
	IsBuyer, err := a.userProvider.IsBuyer(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", sl.Err(err))
			return false, ErrInvalidAppID
		}
		return false, err
	}
	a.log.Info("checking if user is buyer", slog.Int64("userID", userID), slog.Bool("isBuyer", IsBuyer))
	return IsBuyer, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	IsAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", sl.Err(err))
			return false, ErrInvalidAppID
		}
		return false, err
	}
	a.log.Info("checking if user is admin", slog.Int64("userID", userID), slog.Bool("isAdmin", IsAdmin))
	return IsAdmin, nil
}
