package jwt

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/jobboard/internal/jwt"
	"example.com/jobboard/internal/passwords"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	username := "testuser"
	password := "password"
	hashedPassword, _ := passwords.HashPassword(password)
	mock.ExpectQuery("SELECT password FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(hashedPassword))

	reqBody := []byte(`{"username":"` + username + `","password":"` + password + `"}`)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt.LoginHandler(db, w, r)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.NotEmpty(t, rr.Body.String())
}

func TestProtectedRequest_Unauthorized(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := jwt.ProtectedRequest(func(w http.ResponseWriter, r *http.Request) {
		t.Fail()
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	assert.Contains(t, rr.Body.String(), "Missing auth header")
}

// Further tests:
// - TestLoginHandler_FailureInvalidCredentials
// - TestLoginHandler_FailureDBError
// - TestProtectedRequest_Success
// - TestProtectedRequest_InvalidToken
