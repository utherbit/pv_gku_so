package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"pahomov_frolovsky_cson/utilities"
)

var Conn *pgxpool.Pool

func InitPostgresConnection() {
	var host, port, username, password, database, schema string

	utilities.LookupEnv(&host, "POSTGRES_HOST")
	utilities.LookupEnv(&port, "POSTGRES_PORT")
	utilities.LookupEnv(&username, "POSTGRES_USERNAME")
	utilities.LookupEnv(&password, "POSTGRES_PASSWORD")
	utilities.LookupEnv(&database, "POSTGRES_DATABASE_NAME")
	utilities.LookupEnv(&schema, "POSTGRES_SCHEMA", "public")

	cfg, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database))
	utilities.PanicIfErr(err)

	Conn, err = pgxpool.NewWithConfig(context.Background(), cfg)
	utilities.PanicIfErr(err)

	err = Conn.Ping(context.Background())
	utilities.PanicIfErr(err)

	log.Printf("\nConnect to Postgres was successfully")
}
