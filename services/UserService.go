package services

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"messenger-api/config"
	"messenger-api/models"
	"messenger-api/utils"
	"net/http"
	"time"
)

func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные"})
		return
	}

	var exists bool
	if err := config.DB.QueryRow("SELECT exists(SELECT 1 FROM users WHERE email=$1)", req.Email).Scan(&exists); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь уже существует"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки пароля"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	var peerID int
	for {
		peerID = rand.Intn(90000000) + 10000000
		if err := config.DB.QueryRow("SELECT exists(SELECT 1 FROM users WHERE peer_id=$1)", peerID).Scan(&exists); err != nil {
			log.Printf("Error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
			return
		}
		if !exists {
			break
		}
	}

	var userID int
	insertQuery := "INSERT INTO users (email, fullname, password_hash, peer_id) VALUES ($1, $2, $3, $4) RETURNING id"
	if err := config.DB.QueryRow(insertQuery, req.Email, req.FullName, string(hashedPassword), peerID).Scan(&userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации"})
		return
	}

	token, err := utils.GenerateToken(peerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные"})
		return
	}

	var user models.User
	query := "SELECT id, email, fullname, password_hash, peer_id FROM users WHERE email=$1"
	if err := config.DB.QueryRow(query, req.Email).Scan(&user.ID, &user.Email, &user.FullName, &user.PasswordHash, &user.PeerID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	token, err := utils.GenerateToken(user.PeerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
