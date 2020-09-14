package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/abe21412/slack-clone-backend/src/util"
	"github.com/jackc/pgx/v4/pgxpool"
)

var once sync.Once
var connPool *pgxpool.Pool

func Init() {
	once.Do(func() {
		dbCreds := util.GetSecret("slack_db_creds")
		if dbCreds == nil {
			panic("failed to get db creds")
		}
		connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?&pool_max_conns=100", dbCreds["username"], dbCreds["password"], dbCreds["host"], dbCreds["dbname"])
		pool, err := pgxpool.Connect(context.Background(), connString)
		fmt.Println(connString)
		if err != nil {
			fmt.Println(err.Error())
			panic("failed to connect to db")
		}
		connPool = pool
	})
}

func GetPool() *pgxpool.Pool {
	if connPool == nil {
		Init()
	}
	return connPool
}

func Close() {
	connPool.Close()
}
