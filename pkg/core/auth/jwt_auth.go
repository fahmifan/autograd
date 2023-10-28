package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var ErrTokenInvalid = errors.New("token invalid")

type JWTToken string

type Claim struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func CreateTokenExpiry() int64 {
	expireTime := time.Now().Add(8 * time.Hour)
	tokenExpiry := expireTime.UnixNano() / 1_000_000
	return tokenExpiry
}

func GenerateJWTToken(jwtKey string, user AuthUser, expiry int64) (JWTToken, error) {
	claims := &Claim{
		ID:    user.UserID.String(),
		Email: user.Email,
		Role:  user.Role.ToString(),
		Name:  user.Name,
		StandardClaims: jwt.StandardClaims{
			// millisecond
			ExpiresAt: expiry,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return JWTToken(tokenString), nil
}

func ParseJWTClaims(jwtKey string, token JWTToken) (Claim, error) {
	claim := &Claim{}
	tkn, err := jwt.ParseWithClaims(string(token), claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return *claim, err
		}
	}

	if tkn != nil && !tkn.Valid {
		return *claim, ErrTokenInvalid
	}

	return *claim, nil
}

func ParseToken(jwtKey string, token JWTToken) (AuthUser, bool) {
	claims, err := ParseJWTClaims(jwtKey, token)
	if err != nil {
		return AuthUser{}, false
	}

	guid, err := uuid.Parse(claims.ID)
	if err != nil {
		return AuthUser{}, false
	}

	user := AuthUser{
		UserID: guid,
		Email:  claims.Email,
		Role:   Role(claims.Role),
		Name:   claims.Name,
	}

	return user, true
}
