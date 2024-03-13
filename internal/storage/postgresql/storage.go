package postgresql

import (
	"context"
	"errors"
	"tn/internal/domain/models"
	"tn/internal/storage"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(storagePath string) (*Storage, error) {

	db, err := gorm.Open(postgres.Open(storagePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *Storage) SaveUser(ctx context.Context, login string, email string, pwdHash []byte, role int32) (int64, error) {
	user := models.User{
		Login:   login,
		Email:   email,
		PwdHash: pwdHash,
		Role:    role,
	}

	result := s.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return -1, storage.ErrUserExists
		}
		return -1, result.Error
	}
	// добавить вызовы создания в других микросерввисах
	return user.UUID, nil

}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {

	var user models.User
	result := s.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, storage.ErrUserNotFound
		}
		return user, result.Error
	}
	return user, nil
}

// too much duplicate code?
func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.Role == 0, nil
}

func (s *Storage) IsOrginiser(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	return user.Role == 1, nil
}

func (s *Storage) IsDistributor(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.Role == 2, nil
}

func (s *Storage) IsBuyer(ctx context.Context, userID int64) (bool, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.Role == 3, nil
}

// Будет ли рабоать с интефейсом???
func (s *Storage) FindByID(ctx context.Context, userID int64) (models.User, error) {
	var user models.User
	result := s.db.WithContext(ctx).Where("uuid = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return user, storage.ErrUserNotFound
		}
		return user, result.Error
	}
	return user, nil
}

func (s *Storage) Update(ctx context.Context, user models.User) error {
	result := s.db.WithContext(ctx).Save(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return storage.ErrUserNotFound
		}
		return result.Error
	}
	return nil
}
