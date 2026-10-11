package database

import (
	"FumiServer/internal/models"
	"context"
	"log"
	"net/http"
	"os"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var DB *pgxpool.Pool

func InitDatabase() {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Printf("Failed to read DB_URL from .env")
		return
	}
	connString := dbUrl
	var err error

	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Printf("Failed to create db: %v", err)
		return
	}
	log.Printf("Database created successfully")
	initTables(DB)
}

func initTables(db *pgxpool.Pool) {
	query := "CREATE TABLE IF NOT EXISTS users(" +
		"id UUID PRIMARY KEY DEFAULT uuidv7()," +
		"username varchar(16) UNIQUE NOT NULL," +
		"password_hash varchar(70) NOT NULL" +
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

func UserHashPassCompare(db *pgxpool.Pool, r *http.Request, req *models.AuthJson) (uuid.UUID, error) {
	var userId uuid.UUID
	var storedHash string
	user := req.Username
	query := "SELECT  id, password_hash " +
		"FROM users " +
		"WHERE username = $1;"

	err := db.QueryRow(r.Context(), query, user).Scan(&userId, &storedHash)
	if err != nil {
		log.Printf("Failed to find user %v", err)
		return [16]byte{}, err
	}
	compareResult := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password))
	if compareResult != nil {
		log.Printf("Failed to compare password %v", compareResult)
		return [16]byte{}, compareResult
	}
	return userId, compareResult
}
