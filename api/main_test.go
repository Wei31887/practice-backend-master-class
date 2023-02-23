package api

import (
	"lesson/simple-bank/config"
	db "lesson/simple-bank/db/sqlc"
	"lesson/simple-bank/utils"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := config.Config{
		SecreteKey:   utils.RandomString(32),
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}