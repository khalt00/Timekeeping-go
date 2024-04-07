package model

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"timekeeping/db/dbsvc"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
	IsAdmin      bool   `json:"isAdmin"`
}

var jwtSecret = []byte("supersecretkey")

// Claims represents the claims for JWT token
type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

func Login(req LoginRequest) (string, error) {

	user, err := GetUserByUsername(req.Username)
	if user == nil || err != nil {
		log.Println("GetUserByUsername", err)
		return "", errors.New("User not found or some error")
	}

	if user.PasswordHash != req.Password {
		return "", errors.New("password incorrect")
	}
	// // Verify password
	// if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	// 	return
	// }
	// Generate JWT token
	token, err := generateToken(user)
	if err != nil {
		return "", errors.New("can't generateToken")
	}
	return token, nil

}

// GetUserByUsername retrieves a user from the database by username
func GetUserByUsername(username string) (*User, error) {
	conn := dbsvc.GetPostgresConn()

	var user User
	query := "SELECT id, username, password_hash, is_admin FROM users WHERE username = $1"
	err := conn.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("User not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

// generateToken generates JWT token for the user
func generateToken(user *User) (string, error) {
	// Create the claims
	claims := &Claims{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
