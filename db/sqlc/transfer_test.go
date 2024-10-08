package db

import (
	"context"
	"github.com/amer-web/simple-bank/helper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTransfer(t *testing.T) {
	acc := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: acc.ID,
		ToAccountID:   acc2.ID,
		Amount:        helper.RandomInt(1, 1000),
	}
	createTransfer, err := testStore.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.FromAccountID, createTransfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, createTransfer.ToAccountID)
	require.Equal(t, arg.Amount, createTransfer.Amount)

}

//func TestGetAccount(t *testing.T) {
//	account := createRandomAccount(t)
//	account2, err := testQueries.GetAccount(context.Background(), account.ID)
//	require.NoError(t, err)
//	require.NotEmpty(t, account)
//	require.Equal(t, account.Owner, account2.Owner)
//	require.Equal(t, account.Balance, account2.Balance)
//	require.Equal(t, account.Currency, account2.Currency)
//	require.NotZero(t, account2.ID)
//	require.NotZero(t, account2.CreatedAt)
//}
//func TestUpdateAccount(t *testing.T) {
//	account := createRandomAccount(t)
//	arg := UpdateAccountParams{
//		ID:      account.ID,
//		Balance: helper.RandomInt(1, 1000),
//	}
//	accountUpdatae, err := testQueries.UpdateAccount(context.Background(), arg)
//	require.NoError(t, err)
//	require.NotEmpty(t, accountUpdatae)
//	require.Equal(t, arg.ID, accountUpdatae.ID)
//	require.Equal(t, arg.Balance, accountUpdatae.Balance)
//	require.Equal(t, account.Owner, accountUpdatae.Owner)
//
//}
//func TestDeleteAccount(t *testing.T) {
//	account := createRandomAccount(t)
//	err := testQueries.DeleteAccount(context.Background(), account.ID)
//	require.NoError(t, err)
//	account2, err := testQueries.GetAccount(context.Background(), account.ID)
//	require.Error(t, err)
//	require.Equal(t, err, sql.ErrNoRows)
//	require.Empty(t, account2)
//
//}
//func TestListAccounts(t *testing.T) {
//	for i := 0; i < 10; i++ {
//		createRandomAccount(t)
//	}
//	arg := ListAccountsParams{
//		Limit:  5,
//		Offset: 5,
//	}
//	accounts, err := testQueries.ListAccounts(context.Background(), arg)
//	require.NoError(t, err)
//	require.NotEmpty(t, accounts)
//	require.Len(t, accounts, 5)
//
//}
