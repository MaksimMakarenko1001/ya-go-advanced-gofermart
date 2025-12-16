package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Repository struct {
	key            string
	expireInterval time.Duration
}

func New(cfg Config) *Repository {
	return &Repository{
		key:            cfg.Key,
		expireInterval: cfg.ExpireInterval,
	}
}

func (r *Repository) JwtValidate(tokenString string) (ok bool, err error) {
	token, err := jwt.Parse(tokenString,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(r.key), nil
		})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, nil
		}
		return false, err
	}
	return token.Valid, nil
}

func (r *Repository) JwtGenerateToken(userID int64) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(r.expireInterval)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(r.key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *Repository) JwtGetUserID(tokenString string) (int64, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(r.key), nil
	})
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
