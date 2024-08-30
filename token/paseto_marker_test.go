package token

import (
	"github.com/amer-web/simple-bank/helper"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewPasetoMaker(t *testing.T) {
	secret := helper.RandomString(32)
	newPasetoMaker := NewPasetoMaker(secret)
	username := helper.RandomString(3)
	duration := time.Minute
	token, err := newPasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	paylod, err := newPasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, paylod)
	require.Equal(t, username, paylod.Sub)
}
