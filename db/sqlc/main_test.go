package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/spaghetti-lover/bank-system/util"
)
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
    config, err := util.LoadConfig("../..")
    if err != nil {
        log.Fatal("cannot load config file: ", err)
    }
    testDB, err = sql.Open(config.DBDriver, config.DBSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }
	testQueries = New(testDB)
    // Exit with the test result code
    os.Exit(m.Run())
}