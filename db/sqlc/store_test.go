package db

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx_CreateTransferConcurrent(t *testing.T) {
	ctx := context.Background()
	clearTables(ctx, t)

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	n := 5
	amount := int64(10)
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	wg.Wait()
	close(errs)
	close(results)

	existed := make(map[int]bool)

	for err := range errs {
		require.NoError(t, err)
	}

	for result := range results {
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.CreatedAt)

		_, err := testStore.GetTransfer(ctx, transfer.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		toAccount := result.ToAccount

		require.NotEmpty(t, fromAccount)
		require.NotEmpty(t, toAccount)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
}
