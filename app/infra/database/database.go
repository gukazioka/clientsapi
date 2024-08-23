package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gkazioka/clientsapi/app/infra/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setupDatabase(dbPool *pgxpool.Pool) {
	_, err := dbPool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(100), code VARCHAR(30) UNIQUE)")
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
	}

}

func NewDatabase(config config.DatabaseConfig) *pgxpool.Pool {
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Db,
	)

	cfg, _ := pgxpool.ParseConfig(dbUrl)

	cfg.MaxConns = 8
	cfg.MinConns = 2

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to test connection to database.")
	}

	fmt.Fprintln(os.Stdout, "Connection estabilished.")

	setupDatabase(dbPool)

	return dbPool

}
