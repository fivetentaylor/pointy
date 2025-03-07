package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/fivetentaylor/pointy/pkg/utils"
)

func Test_Contains(t *testing.T) {
	tests := []struct {
		name   string
		slice  []string
		values []string
		output bool
	}{
		{"Normal usage", []string{"1", "2"}, []string{"1"}, true},
		{"Not found", []string{"1", "2"}, []string{"3"}, false},
		{"Empty slice", []string{}, []string{"1"}, false},
		{"Empty values", []string{"1", "2"}, []string{}, false},
		{"Empty slice and values", []string{}, []string{}, false},
		{"Multiple values", []string{"1", "2", "3"}, []string{"1", "3"}, true},
		{"Multiple values not found", []string{"1", "2", "3"}, []string{"1", "4"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.Contains(tt.slice, tt.values...)
			assert.Equal(t, tt.output, got)
		})
	}
}
