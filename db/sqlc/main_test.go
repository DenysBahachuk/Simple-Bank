package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	// Import the PostgreSQL driver
	"github.com/DenysBahachuk/Simple_Bank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	cfg, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(cfg.DBdriver, cfg.DBsource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	fmt.Println("connected to db:", cfg.DBsource)

	testQueries = New(testDB)

	// Run tests
	code := m.Run()

	// Clean up and close the database connection
	if err := testDB.Close(); err != nil {
		log.Fatal("failed to close the database:", err)
	}

	os.Exit(code)
}
