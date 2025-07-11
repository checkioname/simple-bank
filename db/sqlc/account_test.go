package db

import (
	"context"
	"fmt"
	"github.com/checkioname/simple-bank/util"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // driver migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	testAccount := CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}

	ctx := context.Background()
	account, err := testQueries.CreateAccount(ctx, testAccount)
	fmt.Println(account)

	require.NoError(t, err)
	require.True(t, account.ID == 1, "Does not match ID")
	require.True(t, account.Balance == testAccount.Balance, "Does not match Balance")
	require.True(t, account.Currency == testAccount.Currency, "Does not match Currency")

	err = testQueries.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)

}

func TestDeleteAccount(t *testing.T) {
	testAccount := CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}

	ctx := context.Background()
	account, err := testQueries.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	err = testQueries.DeleteAccount(ctx, account.ID)
	require.NoError(t, err, "Error deleting account")

	_, err = testQueries.GetAccount(ctx, account.ID)
	require.Error(t, err)
	require.Equal(t, err, pgx.ErrNoRows)
}

func TestGetAccount(t *testing.T) {
	testAccount := CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}
	ctx := context.Background()
	account, err := testQueries.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(ctx, account.ID)
	require.NoError(t, err)
	require.Equal(t, acc.Balance, testAccount.Balance)
	require.Equal(t, acc.Currency, testAccount.Currency)
	require.Equal(t, acc.Owner, testAccount.Owner)

	err = testQueries.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	testAccount := CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}

	ctx := context.Background()
	account, err := testQueries.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	updParams := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}
	acc, err := testQueries.UpdateAccount(ctx, updParams)
	require.NoError(t, err)
	require.Equal(t, updParams.Balance, acc.Balance)

	err = testQueries.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)
}

func TestListAccounts(t *testing.T) {
	listParams := ListAccountsParams{
		Limit: 10,
	}
	ctx := context.Background()
	account, err := testQueries.ListAccounts(ctx, listParams)
	require.NoError(t, err)
	require.Len(t, account, 0)
}
