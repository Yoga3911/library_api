package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func DatabaseConnection() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	dsn := "dev"
	switch dsn {
	case "dev":
		dsn = fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
			os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	case "prod":
		dsn = os.Getenv("DATABASE_URL")
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Println(err.Error())
	}

	config.MaxConns = 20
	config.MinConns = 5

	pg, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Println(err.Error())
	}

	action := 0
	switch action {
	case 1:
		migration(pg)
	case 2:
		rollback(pg)
	}

	return pg
}
