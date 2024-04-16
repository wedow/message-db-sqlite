package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash64(t *testing.T) {
	tests := []struct {
		in  string
		out int64
	}{
		{"123", 144901983852400757},
	}

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			actual, err := hash64Fn(test.in)
			require.NoError(t, err)
			require.Equal(t, test.out, actual)
		})
	}
}
