package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = VerifyPassword(password, hash)
	require.NoError(t, err)

	wrongPassword := "password1010"
	err = VerifyPassword(wrongPassword, hash)
	require.Error(t, err)
	require.Equal(t, err, bcrypt.ErrMismatchedHashAndPassword)
}
