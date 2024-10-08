package db

import (
	"context"
	"github.com/amer-web/simple-bank/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.Username,
		//Balance:  helper.RandomInt(1, 1000),
		Balance:  400,
		Currency: helper.RandomCurrency(),
	}
	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}
func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Balance, account2.Balance)
	require.Equal(t, account.Currency, account2.Currency)
	require.NotZero(t, account2.ID)
	require.NotZero(t, account2.CreatedAt)
}
func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account.ID,
		Balance: helper.RandomInt(1, 1000),
	}
	accountUpdatae, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accountUpdatae)
	require.Equal(t, arg.ID, accountUpdatae.ID)
	require.Equal(t, arg.Balance, accountUpdatae.Balance)
	require.Equal(t, account.Owner, accountUpdatae.Owner)

}
func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testStore.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)
	account2, err := testStore.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	require.Equal(t, err, ErrorRecordNotFound)
	require.Empty(t, account2)

}
func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

}
