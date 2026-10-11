package autorization

import (
	"FumiServer/internal/database"
	"FumiServer/internal/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	var u models.User
	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId := uuid.NewV7()

		hashPass, err := getHashPassword(u.Password)
		if err != nil {
			log.Printf("Failed to create hash password: %v", err)
		}

		udb := models.User{
			Id:       userId,
			Username: u.Username,
			Password: string(hashPass),
		}
		err = database.CreateUser(db, r.Context(), udb)
		if err != nil {
			log.Printf("Failed to record user: %v", err)
		}
	}
}

func LoginHandler(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Printf("Failed to read JWT_SECRET from .env")
		return "", errors.New("failed to read JWT_SECRET from .env")
	}
	var req models.AuthJson
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Failed to decode json: %v", err)
		return "", err
	}
	userId, res := database.UserHashPassCompare(db, r, &req)
	if res != nil {
		return "", res
	}
	if userId == [16]byte{} {
		return "", errors.New("failed to generate token for empty uuid")
	}

	claims := jwt.RegisteredClaims{
		ID:        uuid.New().String(),
		Issuer:    "FumiChat",
		Subject:   userId.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now())}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Failed to sign token: %v", err)
		return "", err
	}
	return signToken, err

}

func getHashPassword(password string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPass), err
}
