package api

import (
	"os"
	"testing"
	"time"

	db "github.com/ifandonlyif-io/ifandonlyif-backend/db/sqlc"
	"github.com/ifandonlyif-io/ifandonlyif-backend/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		AccessTokenSymmetricKey:  util.RandomString(32),
		AccessTokenDuration:      time.Minute,
		RefreshTokenSymmetricKey: util.RandomString(32),
		RefreshTokenDuration:     time.Hour,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	echo.New()

	os.Exit(m.Run())
}
