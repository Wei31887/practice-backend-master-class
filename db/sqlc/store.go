package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store structure for all functions to do queries and transactions
type Store struct {
	*Queries
	db *sql.DB  // for transction db
}

// Create a new Store structure for queries transaction
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}

// execTX: executes the function witn database transaction
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) 
	if err != nil {
		return nil
	}

	qu := New(tx)
	err = fn(qu)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}


// TransferTx perform a money transfer from one account to another account
// Creates the transfer record, account entries, and update accounts' balance in a single transaction.
func (store *Store) TranserTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult


	err := store.execTX(ctx, func(q *Queries) (err error){
		// 1. create the transfer record
		// result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
		// 	FromAccountID: arg.FromAccountID,
		// 	ToAccountID: arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return 
		}

		// 2.1 From entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,	
			Amount: -arg.Amount,
		})
		if err != nil {
			return
		}

		// 2.2 To entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,	
			Amount: arg.Amount,
		})
		if err != nil {
			return
		}

		// TODO: 3. balance two ammount 
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return 
		}
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: account1.ID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return 
		}

		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return 
		}
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: account2.ID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return 
		}


		return nil
	})

	return result, err
}

