package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx_WhatIsBeenTested(t *testing.T) {
	ctx := context.Background()

	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	//errs := make(chan error)
	//results := make(chan TransferTxResult)

	for i := 0; i < 5; i++ {
		//go func() {
		result, err := testStore.TransferTx(ctx, TransferTxParams{
			FromAccountID: acc1.ID,
			ToAccountID:   acc2.ID,
			Amount:        100,
		})
		fmt.Printf("Error transfer: %v", err)
		fmt.Println(result)

		require.NoError(t, err)
		require.NotEmpty(t, result)
		//errs <- err
		//results <- result
		//}()
	}

	//for i := 0; i < 5; i++ {
	//	err := <-errs
	//
	//	result := <-results
	//
	//	transfer := result.Transfer
	//	require.NotEmpty(t, transfer)
	//	require.Equal(t, acc1.ID, transfer.FromAccountID)
	//	require.Equal(t, acc2.ID, transfer.ToAccountID)
	//	require.NotZero(t, transfer.CreatedAt)
	//
	//	_, err = testStore.GetTransfer(ctx, transfer.ID)
	//	require.NoError(t, err)
	//}

}
