package configs

import (
	"context"
	// "fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	// "github.com/joho/godotenv"
)

func DatabaseConnection() *pgxpool.Pool {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	// dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
	// 	os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err.Error())
	}

	config.MaxConns = 20
	config.MinConns = 5

	pg, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Println(err.Error())
	}

	action := 1
	switch action {
	case 1:
		err = migration(pg)
	case 2:
		err = rollback(pg)
	default:
		err = nil
	}

	if err != nil {
		log.Println(err.Error())
	}

	return pg
}
