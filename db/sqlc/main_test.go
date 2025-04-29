package db

import (
<<<<<<< HEAD
	"database/sql"
=======
	"context"
>>>>>>> d4d0e58 (refactor)
	"log"
	"os"
	"testing"

<<<<<<< HEAD
	_ "github.com/lib/pq"
	"github.com/spaghetti-lover/simplebank/util"
)

var testQueries *Queries
var testDB *sql.DB
=======
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spaghetti-lover/simplebank/util"
)

var testStore Store
>>>>>>> d4d0e58 (refactor)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

<<<<<<< HEAD
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
=======
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
>>>>>>> d4d0e58 (refactor)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

<<<<<<< HEAD
	testQueries = New(testDB)

=======
	testStore = NewStore(connPool)
>>>>>>> d4d0e58 (refactor)
	os.Exit(m.Run())
}
