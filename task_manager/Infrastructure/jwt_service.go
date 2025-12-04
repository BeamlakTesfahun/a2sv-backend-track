package Infrastructure

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID, username, role string) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTServiceImpl struct {
	secret []byte
}

func NewJWTService() JWTService {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		sec = "dev_secret_change_me"
	}
	return &JWTServiceImpl{secret: []byte(sec)}
}

func (j *JWTServiceImpl) GenerateToken(userID, username, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTServiceImpl) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}
