package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
    store := NewStore(testDB)

    account1 := createRandomAccount(t)
    account2 := createRandomAccount(t)

    //run n concurrent transfer transaction
    n := 5
    amount := int64(10)
    results := make(chan TransferTxResult)
    errs := make(chan error)
    for i := 0; i < n; i++ {
        go func() {
            result, err := store.TransferTx(context.Background(), TransferTxParams{
                FromAccountID: account1.ID,
                ToAccountID:   account2.ID,
                Amount:        amount,
            })
            errs <- err
            results <- result
        }()
    }

    for i := 0; i < n; i++ {
        err := <-errs
        require.NoError(t, err)

        result := <-results
        require.NotEmpty(t, result)

       //check transfer
	   transfer := result.Transfer
	   require.NotZero(t, transfer.ID)
	   require.Equal(t, account1.ID, transfer.FromAccountID)
	   require.Equal(t, account2.ID, transfer.ToAccountID)
	   require.Equal(t, amount, transfer.Amount)
	   require.NotZero(t, transfer.CreatedAt)
	   require.NotZero(t, transfer.ID)

	   _, err = store.GetTransfer(context.Background(), transfer.ID)
	   require.NoError(t, err)

	   //check entries
	   fromEntry := result.FromEntry
	   require.NotZero(t, fromEntry)
	   require.Equal(t, account1.ID, fromEntry.AccountID)
	   require.Equal(t, -amount, fromEntry.Amount)
	   require.NotZero(t, fromEntry.CreatedAt)
	   require.NotZero(t, fromEntry.ID)
	   require.NotZero(t, fromEntry.CreatedAt)

	   _, err = store.GetEntry(context.Background(), fromEntry.ID)
	   require.NoError(t, err)

	   toEntry := result.ToEntry
	   require.NotZero(t, toEntry)
	   require.Equal(t, account2.ID, toEntry.AccountID)
	   require.Equal(t, amount, toEntry.Amount)
	   require.NotZero(t, toEntry.CreatedAt)
	   require.NotZero(t, toEntry.ID)

	   _, err = store.GetEntry(context.Background(), toEntry.ID)
	   require.NoError(t, err)

	   //check account's balance
    }
}