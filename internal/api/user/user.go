package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/jobboard/internal/jwt"
	"example.com/jobboard/internal/passwords"
	"example.com/jobboard/internal/types"
)

type DecodedToken = types.DecodedToken

type RegisterBody struct {
	Username string
	Password string
}

func ValidateInput(username, password string) bool {
	return true
}

func ReqisterUser(db *sql.DB, w http.ResponseWriter, r*http.Request){
	w.Header().Set("Content-Type", "application/json")

	var body RegisterBody

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !ValidateInput(body.Username, body.Password) {
		http.Error(w, "Username or password does not meet the criteria", http.StatusBadRequest)
		return
	}

	fmt.Printf("The user request value %v\n", body)

	hashPass, hashErr := passwords.HashPassword(body.Password)

	if hashErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", body.Username).Scan(&exists)
	if err != nil || exists > 0 {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
	}

	_, err = db.Exec(
		`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`,
		body.Username, hashPass, time.Now(),
	)
	if err != nil {
		log.Printf("Error inserting user: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User registered successfully")
}

func Iam(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	user, err := jwt.DecodeToken(tokenString) 

	if err != nil {
		log.Printf("Error decoding token: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling user to JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(userJson)
}
