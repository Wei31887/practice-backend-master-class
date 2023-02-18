package db

import (
	"context"
	"testing"

	"lesson/simple-bank/utils"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	hashPd, err := utils.HashedPassword(utils.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashPd,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)

	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangeAt)
	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.Equal(t, user1.FullName, user2.FullName)

	require.NotEmpty(t, user2.CreatedAt)
	require.NotEmpty(t, user2.PasswordChangeAt)
}
