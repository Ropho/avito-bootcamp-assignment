package jwt

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Ropho/avito-bootcamp-assignment/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNotValid = errors.New("not a valid token")
	ErrParse    = errors.New("token parsing error")
	ErrSign     = errors.New("unable to create signature")
)

// JWTService describes an interface for jwt token generation submodule
type Service interface {
	GenerateUserAccessJWT(uuid string, userType string) (string, error)
	VerifyUserAccessToken(tokenString string) (*UserClaims, error)
}

type jwtService struct {
	accessJWTSecret string
	accessTokenExp  int
}

// NewJWTServiceParams describes necessary params to initialize jwtService
type NewJWTServiceParams struct {
	JwtConfig config.JWTServiceConf
}

// NewJWTService is a constructor of JWTService struct
func NewJWTService(p *NewJWTServiceParams) (Service, error) {
	exp, err := strconv.Atoi(p.JwtConfig.AccessJWTExp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration time %s: %w", p.JwtConfig.AccessJWTExp, err)
	}

	return &jwtService{
		accessJWTSecret: p.JwtConfig.AccessJWTSecret,
		accessTokenExp:  exp,
	}, nil
}

func generateJWT(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", ErrSign
	}

	return tokenString, nil
}

func verifyJWT(tokenString, secret string, claims jwt.Claims) error {

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return ErrParse
	}

	if !token.Valid {
		return ErrNotValid
	}

	return nil
}
