package token

import (
	"fmt"
	"github.com/amer-web/simple-bank/helper"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	secret := helper.RandomString(32)
	newJwtMaker := NewJwtMaker(secret)
	username := helper.RandomString(3)
	duration := time.Minute
	token, err := newJwtMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	fmt.Println(token)
	paylod, err := newJwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, paylod)
	require.Equal(t, paylod.Sub, username)
	require.WithinDuration(t, paylod.Exp, paylod.IAT, time.Minute) // difference between two times is within a specific duration
}
