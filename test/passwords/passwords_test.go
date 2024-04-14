package passwords

import (
	"testing"

	"example.com/jobboard/internal/passwords"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hashedPassword, err := passwords.HashPassword(password)
	assert.NoError(t, err, "Hashing the password should not produce an error")

	assert.NotEqual(t, password, hashedPassword, "The hashed password should not be the same as the original password")

	cost, err := bcrypt.Cost([]byte(hashedPassword))
	assert.NoError(t, err, "Fetching the bcrypt cost should not produce an error")
	assert.Equal(t, 14, cost, "The bcrypt cost should be equal to 14")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword123"
	wrongPassword := "wrongPassword123"

	hashedPassword, err := passwords.HashPassword(password)
	assert.NoError(t, err, "Hashing the password should not produce an error")

	isValid := passwords.CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid, "The check should pass for the correct password")

	isValid = passwords.CheckPasswordHash(wrongPassword, hashedPassword)
	assert.False(t, isValid, "The check should fail for the wrong password")
}

// Further tests:
// - Testing `HashPassword` with edge cases, such as an empty string.
// - Testing `CheckPasswordHash` with invalid hash formats.
