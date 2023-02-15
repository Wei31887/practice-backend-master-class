package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store structure for all functions to do queries and transactions
type Store interface {
	Querier
	TranserTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store structure for all functions to do queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB  // for transction db
}

// Create a new Store structure for queries transaction
func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

// execTX: executes the function witn database transaction
func (store *SQLStore) execTX(ctx context.Context, fn func(*Queries) error) error {
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
func (store *SQLStore) TranserTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult


	err := store.execTX(ctx, func(q *Queries) (err error){
		// 1. create the transfer record
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

		// 3. balance two ammount 
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return
			}
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
    account1ID int64,
	amount1 int64,
    account2ID int64,
    amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: account1ID,
		Amount: amount1,
	})
	if err != nil {
		return 
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: account2ID,
		Amount: amount2,
	})
	return
}

