package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := HashedPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestHashPasswordNotDuplicate(t *testing.T) {
	password := RandomString(6)
	hashedPassword1, err := HashedPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	hashedPassword2, err := HashedPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	require.NotEqual(t, hashedPassword1, hashedPassword2)
}


func TestGoodPasswordCompare(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := HashedPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)
}

func TestBadPasswordCompare(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := HashedPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

