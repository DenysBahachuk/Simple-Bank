package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	// Import the PostgreSQL driver
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://admin:adminpassword@localhost:5432/Simple_Bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("unable to connect to db:", err)
	}

	testQueries = New(testDB)

	// Run tests
	code := m.Run()

	// Clean up and close the database connection
	if err := testDB.Close(); err != nil {
		log.Fatal("failed to close the database:", err)
	}

	os.Exit(code)
}
