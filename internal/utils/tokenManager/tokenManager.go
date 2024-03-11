package tokenmanager

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Manager struct {
	signingKey           []byte //super secret
	AccessTokenLifeTime  time.Duration
	RefreshTokenLifeTime time.Duration
}

func NewManager(signingKey []byte) *Manager {
	return &Manager{
		signingKey:           signingKey,
		AccessTokenLifeTime:  time.Duration(1) * time.Minute,
		RefreshTokenLifeTime: time.Duration(30) * time.Minute,
	}
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"resreshToken"`
}

func (m *Manager) IsValidJWT(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return m.signingKey, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			return false, fmt.Errorf("token is not alive")
		}
		return true, nil
	}

	return false, err

}

func (m *Manager) GenerateNewJWTPair(refreshToken string, userID int64) (TokenResponse, error) {
	if valid, err := m.IsValidJWT(refreshToken); err != nil || !valid {
		return TokenResponse{}, err
	}

	newAccessToken, err := m.NewJWT(strconv.FormatInt(userID, 10))
	if err != nil {
		return TokenResponse{}, err
	}

	newRefreshToken, err := m.NewRefreshToken(strconv.FormatInt(userID, 10))
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{newAccessToken, newRefreshToken}, nil
}

func (m *Manager) NewJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.AccessTokenLifeTime).Unix(),
		Subject:   userID,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) NewRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.RefreshTokenLifeTime).Unix(),
		Subject:   userID,
	})

	return token.SignedString([]byte(m.signingKey))

}
