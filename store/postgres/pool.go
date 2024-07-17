package pool

import (
	"context"
	"log"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitConnPool(ctx context.Context) *pgxpool.Pool {
	dbconfig, err := pgxpool.ParseConfig(os.Getenv("CONN_STRING"))
	if err != nil {
		log.Fatalf("Unable to parse db config: %v\n", err)
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}
	dbpool, err := pgxpool.NewWithConfig(ctx, dbconfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return dbpool
}
