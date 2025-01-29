package api

import (
	"os"
	"testing"
	"time"

	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func newTestServer(store db.Store) (*Server, error) {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	config := utils.Config{
		TokenSymmetricKey: utils.RandomString(32),
		TokenDuration:     time.Minute,
	}

	return NewServer(store, logger, config)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
