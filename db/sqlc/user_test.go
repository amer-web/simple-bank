package db

import (
	"context"
	"github.com/amer-web/simple-bank/helper"
	"github.com/jackc/pgx/v5/pgtype"
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
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	return user
}
func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	getUser, err := testStore.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, getUser)
	require.Equal(t, user.Username, getUser.Username)
	require.Equal(t, user.Email, getUser.Email)
	require.Equal(t, user.FullName, getUser.FullName)
}
func TestUpdateUserOnlyEmail(t *testing.T) {
	user := createRandomUser(t)
	newEmail := helper.RandomEmail()
	updated, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: user.Username,
		Email: pgtype.Text{
			Valid:  true,
			String: newEmail,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, updated)
	require.Equal(t, user.Username, updated.Username)
	require.Equal(t, newEmail, updated.Email)
	require.Equal(t, user.FullName, updated.FullName)

}
