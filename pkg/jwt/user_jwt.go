package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UUID     string
	UserType string
}

func (j *jwtService) GenerateUserAccessJWT(uuid string, userType string) (string, error) {

	return generateJWT(&UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(j.accessTokenExp))),
		},
		UUID:     uuid,
		UserType: userType,
	}, j.accessJWTSecret)
}

func (j *jwtService) VerifyUserAccessToken(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	err := verifyJWT(tokenString, j.accessJWTSecret, claims)

	return claims, err
}
