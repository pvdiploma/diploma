package tokenmanager

import (
	"fmt"
	"strconv"
	"time"
	"tn/internal/domain/models"

	"github.com/dgrijalva/jwt-go"
)

type TokenManager struct {
	signingKey           []byte
	AccessTokenLifeTime  time.Duration
	RefreshTokenLifeTime time.Duration
}

func NewManager(signingKey []byte) *TokenManager {
	return &TokenManager{
		signingKey:           signingKey,
		AccessTokenLifeTime:  time.Duration(1) * time.Minute,
		RefreshTokenLifeTime: time.Duration(30) * time.Minute,
	}
}

type Claims struct {
	UserID int64  `json:"userID"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
	jwt.StandardClaims
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"resreshToken"`
}

func (m *TokenManager) IsValidJWT(tokenString string) (bool, error) {
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

func (m *TokenManager) GenerateNewJWTPair(user models.User) (TokenResponse, error) {
	newAccessToken, err := m.NewJWT(user, m.AccessTokenLifeTime)
	if err != nil {
		return TokenResponse{}, err
	}

	newRefreshToken, err := m.NewJWT(user, m.RefreshTokenLifeTime)
	if err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{newAccessToken, newRefreshToken}, nil
}

func (m *TokenManager) NewJWT(user models.User, duration time.Duration) (string, error) {
	claims := Claims{
		UserID: user.UUID,
		Email:  user.Email,
		Role:   strconv.Itoa(int(user.Role)),
		Exp:    time.Now().Add(duration).Unix(),

		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// app sercet or common tokenManager secret???
	return token.SignedString([]byte(m.signingKey))
}
