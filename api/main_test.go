package api

import (
	"os"
	"testing"

	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func newTestServer(store db.Store) *Server {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	return NewServer(store, logger)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
