package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransfer_CreateTransfer(t *testing.T) {
	ctx := context.Background()

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	transfer := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        acc1.Balance,
	}

	result, err := testStore.CreateTransfer(ctx, transfer)
	require.NoError(t, err)
	require.NotEmpty(t, result.ID)
	require.Equal(t, transfer.FromAccountID, result.FromAccountID)
	require.Equal(t, transfer.ToAccountID, result.ToAccountID)
	require.Equal(t, transfer.Amount, result.Amount)

	clearTables(ctx, t)
}

func TestTransferTx_ListTransfer(t *testing.T) {
	ctx := context.Background()
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	p := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Limit:         10,
		Offset:        0,
	}
	result, err := testQueries.ListTransfers(ctx, p)
	require.NoError(t, err)
	require.Empty(t, result) //empty cause no transfer was created

	createP := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        acc1.Balance,
	}
	_, err = testQueries.CreateTransfer(ctx, createP)
	require.NoError(t, err)

	result, err = testQueries.ListTransfers(ctx, p)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Len(t, result, 1)

	clearTables(ctx, t)
}

func TestQueries_GetTransfer(t *testing.T) {
	ctx := context.Background()

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	result, err := testQueries.GetTransfer(ctx, -1)
	require.Error(t, err)    //no rows in result set
	require.Empty(t, result) //empty cause no transfer was created

	createP := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        acc1.Balance,
	}
	transfer, err := testStore.CreateTransfer(ctx, createP)
	require.NoError(t, err)

	result, err = testStore.GetTransfer(ctx, transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.Equal(t, acc1.ID, transfer.FromAccountID)
	require.Equal(t, acc2.ID, transfer.ToAccountID)

	clearTables(ctx, t)
}
