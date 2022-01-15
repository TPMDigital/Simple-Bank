package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTransferToAccount(t *testing.T, toAccount Account) Transfer {
	fromAccount := createRandomAccount(t)

	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        fromAccount.Balance,
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	toAccount := createRandomAccount(t)
	createRandomTransferToAccount(t, toAccount)
}

func TestGetTransfer(t *testing.T) {
	toAccount := createRandomAccount(t)
	transfer1 := createRandomTransferToAccount(t, toAccount)
	transfer2, err2 := testQueries.GetTransfer(context.Background(), transfer1.ID)
	//
	require.NoError(t, err2)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfer(t *testing.T) {
	toAccount := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransferToAccount(t, toAccount)
	}

	arg := ListTransfersParams{
		Limit:       5,
		Offset:      5,
		ToAccountID: toAccount.ID,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}
}
