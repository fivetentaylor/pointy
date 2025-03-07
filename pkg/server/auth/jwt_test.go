package auth_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/fivetentaylor/pointy/pkg/models"
	"github.com/fivetentaylor/pointy/pkg/server/auth"
)

func TestGenerateUserToken(t *testing.T) {
	type args struct {
		user *models.User
	}
	tests := []struct {
		secretkey string
		name      string
		args      args
		want      jwt.MapClaims
		wantErr   bool
	}{
		{
			name:      "test generate user token",
			secretkey: "tacotacotacotaco",
			args: args{
				user: &models.User{
					ID:    "test",
					Email: "test@t.com",
				},
			},
			want: map[string]interface{}{
				"uid": "test",
				"e":   "test@t.com",
				"exp": float64(time.Now().Add(time.Hour * 24 * 14).Unix()),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := auth.NewJWT(tt.secretkey)
			token, err := e.GenerateUserToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			payload, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(tt.secretkey), nil
			})
			assert.NoError(t, err)

			got := payload.Claims.(jwt.MapClaims)

			fmt.Println("got:", got, "token:", token)

			assert.Equal(t, len(tt.want), len(got), "len(got) != len(want)")

			for k, v := range tt.want {
				if k == "exp" {
					assert.True(t, math.Abs(float64(got[k].(float64)-v.(float64))) <= 10)
					continue
				}

				assert.Equal(t, v, got[k])
			}
		})
	}
}

func TestParseUserToken(t *testing.T) {
	validUserToken, err := auth.NewJWT("tacotacotacotaco").GenerateUserToken(&models.User{
		ID:    "test",
		Email: "test@t.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	expiredUserToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiJ0ZXN0IiwiZSI6InRlc3RAdC5jb20iLCJleHAiOjE2OTY3OTMyMTN9.JDZDLb3-Hrhe4gS32vjXYVudcP98ychXLZMY3H_du7Y"

	type args struct {
		user *models.User
	}
	tests := []struct {
		name      string
		secretkey string
		token     string
		want      *models.UserClaims
		wantErr   bool
	}{
		{
			name:      "test parse broken token",
			secretkey: "tacotacotacotaco",
			token:     "broken",
			wantErr:   true,
		},
		{
			name:      "test parse user token",
			secretkey: "tacotacotacotaco",
			token:     validUserToken,
			want: &models.UserClaims{
				Id:    "test",
				Email: "test@t.com",
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
				},
			},
			wantErr: false,
		},
		{
			name:      "test parse expired user token",
			secretkey: "tacotacotacotaco",
			token:     expiredUserToken,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := auth.NewJWT(tt.secretkey)

			got, err := e.ParseUserToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUserToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			assert.Equal(t, tt.want.Id, got.Id)
			assert.Equal(t, tt.want.Email, got.Email)
			assert.True(t, math.Abs(float64(got.StandardClaims.ExpiresAt-tt.want.StandardClaims.ExpiresAt)) <= 10)
		})
	}
}
