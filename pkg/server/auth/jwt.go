package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/golang-jwt/jwt"
	"github.com/teamreviso/code/pkg/constants"
	"github.com/teamreviso/code/pkg/models"
)

var UserTokenExpiresIn = time.Hour * 24 * 14 // 14 days

type JWTEngine struct {
	secretkey string
}

func NewJWT(secretkey string) *JWTEngine {
	return &JWTEngine{secretkey: secretkey}
}

func (e *JWTEngine) GenerateUserToken(user *models.User) (string, error) {
	claims := &models.UserClaims{
		Id:    user.ID,
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(UserTokenExpiresIn).Unix(),
		},
	}

	if user.Admin {
		log.Info("admin signin")
		claims.Admin = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(e.secretkey))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tokenString, nil
}

func (e *JWTEngine) ParseUserToken(tokenStr string) (*models.UserClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(e.secretkey), nil
	})
	if err != nil {
		return nil, err
	}

	if claimMap, ok := token.Claims.(jwt.MapClaims); ok {
		if !token.Valid {
			return nil, fmt.Errorf("token is not valid")
		}

		id, ok := claimMap["uid"].(string)
		if !ok {
			return nil, fmt.Errorf("token is not valid")
		}

		claims := &models.UserClaims{
			Id:    id,
			Email: claimMap["e"].(string),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: int64(claimMap["exp"].(float64)),
			},
		}

		if admin, ok := claimMap["a"].(bool); ok {
			claims.Admin = admin
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (e *JWTEngine) Attach(ctx context.Context) context.Context {
	return context.WithValue(ctx, constants.JWTContextKey, e)
}

func JWT(ctx context.Context) *JWTEngine {
	obj, ok := ctx.Value(constants.JWTContextKey).(*JWTEngine)
	if !ok {
		panic(fmt.Errorf("jwt engine not found in context"))
	}

	return obj
}
