package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"sync"
	"testing"
)

func TestCreateTrans(t *testing.T) {
	store := NewStore(dbTest)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	log.Println("first  ", acc1.Balance, acc2.Balance)
	n := 10
	chanTrans := make(chan ResultTransfer, n)
	chanErr := make(chan error, n)
	amount := int64(10)
	var waitGroup sync.WaitGroup

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("transfer-%d", i)
		waitGroup.Add(1)
		go func(name string) {
			ctx := context.WithValue(context.Background(), "amount", name)
			defer waitGroup.Done()
			transfer, err := store.TransferTx(ctx, ArrgTransfer{
				FromAcc: acc1.ID,
				ToAcc:   acc2.ID,
				Amount:  amount,
			})
			if err != nil {
				chanErr <- err
			}
			chanTrans <- transfer
		}(name)
	}
	waitGroup.Wait()
	close(chanErr)
	close(chanTrans)
	for err := range chanErr {
		log.Println("transfer err:", err.Error())
		require.NoError(t, err)
	}
	for result := range chanTrans {

		require.NotEmpty(t, result.Transfer)
		require.Equal(t, amount, result.Transfer.Amount)
		// check from entry
		require.Equal(t, acc1.ID, result.EntryFrom.AccountID)
		require.NotEmpty(t, result.EntryFrom)
		require.Equal(t, -amount, result.EntryFrom.Amount)
		require.NotZero(t, result.EntryFrom.Amount)
		require.NotZero(t, result.EntryFrom.CreatedAt)

		// check to entry
		require.Equal(t, acc2.ID, result.EntryTo.AccountID)
		require.NotEmpty(t, result.EntryTo)
		require.Equal(t, amount, result.EntryTo.Amount)
		require.NotZero(t, result.EntryTo.Amount)
		require.NotZero(t, result.EntryTo.CreatedAt)

		// check accounts and balance
		fromAcc := result.FromAcc
		require.Equal(t, acc1.ID, fromAcc.ID)
		require.NotEmpty(t, fromAcc)
		toAcc := result.ToAcc
		require.NotEmpty(t, toAcc)
		require.Equal(t, acc2.ID, toAcc.ID)
	}
	getAcc1, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAcc1)
	getAcc2, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, getAcc2)
	log.Println("after update ", getAcc1.Balance, getAcc2.Balance)

	require.Equal(t, acc1.Balance-(int64(n)*amount), getAcc1.Balance)
	require.Equal(t, acc2.Balance+(int64(n)*amount), getAcc2.Balance)
}
