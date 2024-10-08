package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Store interface {
	Querier
	TransferTx(context context.Context, arrg ArrgTransfer) (ResultTransfer, error)
}
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {

	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

func (s *SQLStore) execTx(context context.Context, fn func(q *Queries) error) error {
	tx, err := s.connPool.Begin(context)
	if err != nil {
		log.Fatal(err.Error())
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if err := tx.Rollback(context); err != nil {
			log.Println("Rollback failed:", err)
		}
		return err
	}
	return tx.Commit(context)
}

type ArrgTransfer struct {
	FromAcc int64 `json:"from_acc"`
	ToAcc   int64 `json:"to_acc"`
	Amount  int64 `json:"amount"`
}
type ResultTransfer struct {
	FromAcc   Account  `json:"from_acc"`
	ToAcc     Account  `json:"to_acc"`
	Transfer  Transfer `json:"transfer"`
	EntryFrom Entry    `json:"entry_from"`
	EntryTo   Entry    `json:"entry_to"`
}

// add tra , add entr , update balance
func (s *SQLStore) TransferTx(context context.Context, arrg ArrgTransfer) (ResultTransfer, error) {
	var result ResultTransfer
	err := s.execTx(context, func(q *Queries) error {
		var err error

		CreatTrans := CreateTransferParams{
			FromAccountID: arrg.FromAcc,
			ToAccountID:   arrg.ToAcc,
			Amount:        arrg.Amount,
		}

		result.Transfer, err = q.CreateTransfer(context, CreatTrans)

		if err != nil {
			return err
		}
		CreateEntryParamsFrom := CreateEntryParams{
			Amount:    -arrg.Amount,
			AccountID: arrg.FromAcc,
		}
		result.EntryFrom, err = q.CreateEntry(context, CreateEntryParamsFrom)
		if err != nil {
			return err
		}
		CreateEntryParamsTo := CreateEntryParams{
			Amount:    arrg.Amount,
			AccountID: arrg.ToAcc,
		}

		result.EntryTo, err = q.CreateEntry(context, CreateEntryParamsTo)
		if err != nil {
			return err
		}

		//acc1, err := q.GetAccountForUpdate(context, arrg.FromAcc)
		//if err != nil {
		//}
		//
		//result.FromAcc, err = q.UpdateAccount(context, UpdateAccountParams{
		//	ID:      arrg.FromAcc,
		//	Balance: acc1.Balance - arrg.Amount,
		//})
		//
		//acc2, err := q.GetAccountForUpdate(context, arrg.ToAcc)
		//if err != nil {
		//	log.Println(err.Error())
		//}
		//result.ToAcc, err = q.UpdateAccount(context, UpdateAccountParams{
		//	ID:      arrg.ToAcc,
		//	Balance: acc2.Balance + arrg.Amount,
		//})
		result.FromAcc, err = q.UpdateAccountBalance(context, UpdateAccountBalanceParams{
			ID:     arrg.FromAcc,
			Amount: -arrg.Amount,
		})
		if err != nil {
			log.Println(err.Error())
		}
		result.ToAcc, err = q.UpdateAccountBalance(context, UpdateAccountBalanceParams{
			ID:     arrg.ToAcc,
			Amount: arrg.Amount,
		})
		if err != nil {
			log.Println(err.Error())
		}
		return nil
	})
	return result, err
}
