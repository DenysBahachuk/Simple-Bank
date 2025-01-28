package db

import (
	"context"
	"testing"
	"time"

	"github.com/DenysBahachuk/Simple_Bank/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPass, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)

	payload := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), payload)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, payload.Username, user.Username)
	require.Equal(t, payload.HashedPassword, user.HashedPassword)
	require.Equal(t, payload.FullName, user.FullName)
	require.Equal(t, payload.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}
