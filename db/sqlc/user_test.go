package db

import (
	"context"
	"github.com/amer-web/simple-bank/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := helper.HashPassword(helper.RandomString(4))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username: helper.RandomString(4),
		FullName: helper.RandomString(4) + " " + helper.RandomString(5),
		Email:    helper.RandomEmail(),
		Password: hashPassword,
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	getUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, getUser)
	require.Equal(t, user.Username, getUser.Username)
	require.Equal(t, user.Email, getUser.Email)
	require.Equal(t, user.FullName, getUser.FullName)
}
