package env

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/fivetentaylor/pointy/pkg/models"
)

func TestContextValue(t *testing.T) {
	ctx := UserClaimCtx(context.Background(), &models.UserClaims{
		Id:    "c1a2b3c4-d5e6-47f8-9f01-123456789abc",
		Email: "taylor@revi.so",
	})

	val, err := UserClaim(ctx)
	assert.NoError(t, err)
	t.Logf("Value: %+v", val)
}

func TestContextNoValue(t *testing.T) {
	ctx := context.Background()
	_, err := UserClaim(ctx)
	if err != nil {
		contains_type := strings.Contains(err.Error(), "with type")
		assert.True(t, contains_type, "error message should contain type info")
		t.Log(err)
	}
}
