package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	testStore := NewStore(testDB)
	
	// test situation: two account make the several transfers
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// run n concorrent transfer transcations
	n := 5
	amount := int64(10)

	errs := make(chan error) 
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TranserTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}

	expored := make(map[int]bool, 0)
	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

		result := <- results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check: from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check: to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
	
		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check: balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		
		// check: balance
		fmt.Println("TX: ", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance 

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		transferTime := int(diff1 / amount)
		require.True(t, transferTime >= 1 && transferTime <= n)
		require.False(t, expored[transferTime])
		
		expored[transferTime] = true
	}

	// check the final updated account
	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)

	fmt.Println("UPDATED: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, updateAccount1.Balance, account1.Balance - int64(n)*amount)
	require.Equal(t, updateAccount2.Balance, account2.Balance + int64(n)*amount)
}

func TestTransferTxDeadLock(t *testing.T) {
	testStore := NewStore(testDB)
	
	// test situation: two account make the several transfers
	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// run n concorrent transfer transcations
	n := 20
	amount := int64(10)

	errs := make(chan error) 
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID

        if i % 2 == 1 {
			fromAccountId = account2.ID
            toAccountId = account1.ID
		}
		go func() {
			_, err := testStore.TranserTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountId,
				ToAccountID: toAccountId,
				Amount: amount,
			})
			errs <- err
		}()
	}

	// check the final updated account
	updateAccount1, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)

	updateAccount2, err := testStore.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)

	fmt.Println("UPDATED: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, updateAccount1.Balance, account1.Balance)
	require.Equal(t, updateAccount2.Balance, account2.Balance)
}