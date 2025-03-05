package models

import (
	"github.com/golang-jwt/jwt"
)

type UserClaims struct {
	Id    string `json:"uid"`
	Email string `json:"e"`
	Admin bool   `json:"a,omitempty"`

	jwt.StandardClaims
}

func (c *UserClaims) Valid() error {
	return c.StandardClaims.Valid()
}
