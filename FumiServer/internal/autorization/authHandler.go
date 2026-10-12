package autorization

import (
	"FumiServer/internal/database"
	"FumiServer/internal/models"
	"encoding/json"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		udb := models.User{
			Id:       userId,
			Username: u.Username,
			Password: string(hashPass),
		}
		err = database.CreateUser(db, r.Context(), udb)
		if err != nil {
			log.Printf("Failed to record user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func LoginHandler(db *pgxpool.Pool, w http.ResponseWriter, r *http.Request) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Printf("Failed to read JWT_SECRET from .env")
		return
	}
	var req models.AuthJson
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Failed to decode json: %v", err)

	}
	userId, err := database.UserHashPassCompare(db, r, &req)
	if err != nil {
		return
	}
	if userId == [16]byte{} {
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		JwtTOKEN: signToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Ошибка при отправке JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getHashPassword(password string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPass), err
}
