package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gocheck/internal/database"
	"gocheck/internal/models"
	"gocheck/internal/auth"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	_, err := database.DB.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", input.Username, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Użytkownik już istnieje lub błąd bazy"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Zarejestrowano pomyślnie"})
}

func Login(c *gin.Context) {
	var input models.User
	var dbUser models.User
	var passwordHash string

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username=$1", input.Username).Scan(&dbUser.ID, &dbUser.Username, &passwordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nieprawidłowy login lub hasło"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Nieprawidłowy login lub hasło"})
		return
	}

	token, _ := auth.GenerateToken(dbUser.ID)
	c.JSON(http.StatusOK, gin.H{"token": token})
}