package helper

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	s := RandomString(6)
	has, err := HashPassword(s)
	require.NoError(t, err)
	require.NotEmpty(t, has)
	err = CheckPasswordHash(s, has)
	require.NoError(t, err)
	err = CheckPasswordHash("wrongpassword", has)
	require.Error(t, err)
}
