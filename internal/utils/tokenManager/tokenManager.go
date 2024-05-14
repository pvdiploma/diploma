package tokenmanager

import (
	"fmt"
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
		AccessTokenLifeTime:  time.Duration(30) * time.Minute,
		RefreshTokenLifeTime: time.Duration(30) * time.Minute,
	}
}

type Claims struct {
	UserID int64  `json:"userID"`
	Email  string `json:"email"`
	Role   int32  `json:"role"`
	Exp    int64  `json:"exp"`
	jwt.StandardClaims
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"resreshToken"`
}

// refactror???
func (m *TokenManager) ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return m.signingKey, nil
	})
}

func (m *TokenManager) IsOrganizer(tokenStr string) (bool, int64) {
	token, err := m.ParseToken(tokenStr)

	if err != nil {
		return false, -1
	}
	fmt.Println(token.Claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			return false, -1
		}
		if int32(claims["role"].(float64)) != 1 {
			return false, -1
		}
		return true, int64(claims["userID"].(float64))
	}

	return false, -1
}

// TODO: refactor later
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
		UserID: user.Id,
		Email:  user.Email,
		Role:   user.Role,
		Exp:    time.Now().Add(duration).Unix(),

		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.signingKey))
}
