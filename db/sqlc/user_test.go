package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/ifandonlyif-io/ifandonlyif-backend/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams{
		FullName:     sql.NullString{util.RandomOwner(), true},
		Wallet:       sql.NullString{util.RandomWalletAddress(), true},
		CountryCode:  sql.NullString{util.RandomCountry(), true},
		EmailAddress: sql.NullString{util.RandomEmail(), true},
		TwitterName:  sql.NullString{util.RandomOwner(), true},
		ImageUri:     sql.NullString{"https://img.seadn.io/files/2ed3306fc4808ae7bc0b75802ea78c95.png?fit=max", true},
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Wallet, user.Wallet)
	require.Equal(t, arg.CountryCode, user.CountryCode)
	require.Equal(t, arg.TwitterName, user.TwitterName)
	require.Equal(t, arg.ImageUri, user.ImageUri)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)
}

func TestGetUserByWallet(t *testing.T) {

	user, err := testQueries.GetUserByWalletAddress(context.Background(), sql.NullString{util.RandomWalletAddress(), true})
	require.NoError(t, err)
	require.NotEmpty(t, user)
}
