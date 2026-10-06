package autorization

import (
	"FumiServer/internal/database"
	"FumiServer/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"uuid"

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

		hashPass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to create hash password: %e", err)
		}

		udb := models.User{
			Id:       userId,
			Username: u.Username,
			Password: string(hashPass),
		}
		err = database.CreateUser(db, r.Context(), udb)
		if err != nil {
			log.Printf("Failed to record user: %e", err)
		}
		return //future JWT

	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func getHashPassword(password string) (string, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPass), err
}

type authJson struct {
	login    string `json:"login"`
	password string `json:"password"`
}
