package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hash, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hash)

	err = CheckPassword(password, hash)
	require.NoError(t, err)

	wongPassword := RandomString(6)

	err = CheckPassword(wongPassword, hash)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
