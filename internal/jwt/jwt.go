package jwt

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/jobboard/internal/passwords"
	"example.com/jobboard/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

type User = types.User
type DecodedToken = types.DecodedToken

func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
	}

	var storedHash string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", u.Username).Scan(&storedHash)
	if err != nil {
			if err == sql.ErrNoRows {
					http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
					http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
	}

	if !passwords.CheckPasswordHash(u.Password, storedHash) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
	}

	tokenString, err := createToken(u.Username)
	if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, tokenString)
}


type Callback func(http.ResponseWriter, *http.Request)

func ProtectedRequest(callback Callback) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing auth header\n")
			return
		}

		tokenString = tokenString[len("Bearer "):]

		err := verifyToken(tokenString)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Invalid token\n")
			return
		}

		callback(w, r)
	}
}

func DecodeToken(tokenString string) (DecodedToken, error) {
	token, err := jwt.ParseWithClaims(tokenString[len("Bearer "):], &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
	})

	if err != nil {
			return DecodedToken{}, err
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
			return DecodedToken{}, fmt.Errorf("invalid token")
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
			return DecodedToken{}, fmt.Errorf("unable to extract user from token")
	}

	exp, ok := (*claims)["exp"].(float64) 
	if !ok {
			return DecodedToken{}, fmt.Errorf("unable to extract expiration from token")
	}

	return DecodedToken{
			Username:   username,
			Exp: int64(exp),
	}, nil
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func verifyToken(tokenString string) error{
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
 })

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// func refreshToken(decodedToken DecodedToken) (string, error) {
// 	expirationTime := time.Now().Add(30 * time.Minute)
// 	claims := &jwt.MapClaims{
// 		"username": decodedToken.Username,
// 		"exp":      expirationTime.Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(secretKey)
// }
