package database

import (
	"FumiServer/internal/models"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDatabase() {
	connString := "postgres://postgres:FumiChat@localhost:5432/FumiDB"
	var err error

	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("Failed to create db: %e", err)
		return
	}
	log.Printf("Database created successfully")
	initTables(DB)
}

func initTables(db *pgxpool.Pool) {
	query := "CREATE TABLE IF NOT EXISTS users(" +
		"id UUID PRIMARY KEY DEFAULT uuidv7()," +
		"username varchar(16) UNIQUE NOT NULL," +
		"password_hash varchar(60) NOT NULL" +
		")"

	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Failed to create tables: %e", err)
	}
	log.Println("Tables created successfully")
}

func CreateUser(db *pgxpool.Pool, ctx context.Context, u models.User) error {
	query := "INSERT INTO users " +
		"(id, username, password_hash)" +
		"VALUES($1,$2,$3)"

	_, err := db.Exec(ctx, query, u.Id, u.Username, u.Password)
	if err != nil {
		log.Printf("Failed to record user")
	}
	log.Println("User has been recorded")
	return err
}
