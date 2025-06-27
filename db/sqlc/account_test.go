package db

import (
	"context"
	"github.com/checkioname/simple-bank/util"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // driver migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq" // Postgres driver
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var db *Queries
var t testing.T
var teardown func()

func TestMain(m *testing.M) {
	db, teardown = StartDBWithTestContainer(&t)
	defer teardown()

	code := m.Run()

	os.Exit(code)
}

func TestCreateAccount(t *testing.T) {
	testAccount := CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}

	ctx := context.Background()
	account, err := db.CreateAccount(ctx, testAccount)

	require.NoError(t, err)
	require.True(t, account.ID == 1, "Does not match ID")
	require.True(t, account.Balance == testAccount.Balance, "Does not match Balance")
	require.True(t, account.Currency == testAccount.Currency, "Does not match Currency")

	_, err = db.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)

}

func TestDeleteAccount(t *testing.T) {
	testAccount := CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}

	ctx := context.Background()
	account, err := db.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	cmdTag, err := db.DeleteAccount(ctx, account.ID)

	require.True(t, cmdTag.RowsAffected() == 1, "No account deleted")
	require.NoError(t, err, "Error deleting account")

	_, err = db.GetAccount(ctx, account.ID)
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
	account, err := db.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	acc, err := db.GetAccount(ctx, account.ID)
	require.NoError(t, err)
	require.Equal(t, acc.Balance, testAccount.Balance)
	require.Equal(t, acc.Currency, testAccount.Currency)
	require.Equal(t, acc.Owner, testAccount.Owner)

	_, err = db.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)
}

func TestUpdateAccount(t *testing.T) {
	testAccount := CreateAccountParams{
		util.RandomOwner(),
		util.RandomMoney(),
		util.RandomCurrency(),
	}

	ctx := context.Background()
	account, err := db.CreateAccount(ctx, testAccount)
	require.NoError(t, err)

	updParams := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}
	acc, err := db.UpdateAccount(ctx, updParams)
	require.NoError(t, err)
	require.Equal(t, updParams.Balance, acc.Balance)

	_, err = db.DeleteAccount(ctx, account.ID)
	require.NoError(t, err)
}

func TestListAccounts(t *testing.T) {
	listParams := ListAccountsParams{
		Limit: 10,
	}
	ctx := context.Background()
	account, err := db.ListAccounts(ctx, listParams)
	require.NoError(t, err)
	require.Len(t, account, 0)
}
