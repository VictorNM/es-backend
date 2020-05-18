package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/victornm/es-backend/pkg/errorutil"
)

type JWTService interface {
	ParseToken(tokenString string) (*UserAuthDTO, error)
	generateToken(u *User) (string, error)
}

type jwtService struct {
	secret  string
	expired time.Duration
}

func NewJWTService(secret string, expiredHours int) JWTService {
	return &jwtService{
		secret:  secret,
		expired: time.Duration(expiredHours) * time.Hour,
	}
}

func (s *jwtService) generateToken(u *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.expired).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "auth.service",
		},
		UserAuthDTO: &UserAuthDTO{UserID: u.ID},
	})

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errorutil.Wrap(ErrUnknown, err)
	}

	return tokenString, nil
}

func (s *jwtService) ParseToken(tokenString string) (*UserAuthDTO, error) {
	var claims jwtClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errorutil.Wrap(ErrNotAuthenticated, err)
	}

	return claims.UserAuthDTO, nil
}
