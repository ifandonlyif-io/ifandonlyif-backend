package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		FullName:      "haha",
		WalletAddress: "0xAdC41d839b7fC82Fb76bF57fAB7cdDf83bFa68aC",
		CountryCode:   "TW",
		EmailAddress:  "Test@gmail.com",
		TwitterName:   "happy123",
		ImageUri:      "https://img.seadn.io/files/2ed3306fc4808ae7bc0b75802ea78c95.png?fit=max",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.WalletAddress, user.WalletAddress)
	require.Equal(t, arg.CountryCode, user.CountryCode)
	require.Equal(t, arg.TwitterName, user.TwitterName)
	require.Equal(t, arg.ImageUri, user.ImageUri)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
}
