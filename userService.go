package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func RegisterUserHandler(c *gin.Context) {
	var newUser User
	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := UsernameExists(newUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser.Password = hashedPassword

	err = SaveUserToRedis(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})

}

func UserProfileHandler(c *gin.Context) {
	user, _ := c.Get("user")
	authUser := user.(*User)

	c.JSON(http.StatusOK, gin.H{
		"username": authUser.Username,
	})
}

func UpdateUsernameHandler(c *gin.Context) {
	user, _ := c.Get("user")
	authUser := user.(*User)

	var updatedUsername struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBind(&updatedUsername); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedUsername.Username != "" && updatedUsername.Username != authUser.Username {
		exists, err := UsernameExists(updatedUsername.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
	}

	if len(updatedUsername.Username) < 3 || len(updatedUsername.Username) > 20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username length must be between 3 and 20"})
		return
	}

	userID := strings.Split(authUser.Username, ":")[1]

	if err := UpdateUsernameInRedis(userID, updatedUsername.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	authUser.Username = updatedUsername.Username

	c.JSON(http.StatusOK, gin.H{"message": "Username updated", "username": authUser.Username})

}

func LoginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := FindUserByUsername(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if dbUser == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
		return
	}

	if !ComparePasswords(dbUser.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": dbUser.Username,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func UsernameExists(username string) (bool, error) {
	keys, err := rdb.Keys(ctx, "user:*").Result()
	if err != nil {
		return false, err
	}

	for _, key := range keys {
		user, err := rdb.HGet(ctx, key, "username").Result()
		if err != nil {
			return false, err
		}
		if user == username {
			return true, nil
		}
	}

	return false, nil
}

func SaveUserToRedis(user User) error {
	userID, err := rdb.Incr(ctx, "user_id_counter").Result()
	if err != nil {
		return err
	}

	_, err = rdb.HSet(ctx, fmt.Sprintf("user:%d", userID), "username", user.Username, "password", user.Password).Result()
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func UpdateUsernameInRedis(userID string, newUsername string) error {
	userKey := fmt.Sprintf("user:%d", userID)

	exists, err := rdb.Exists(ctx, userKey).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return fmt.Errorf("user with ID %s does not exist", userID)
	}

	if err := rdb.HSet(ctx, userKey, "username", newUsername).Err(); err != nil {
		return err
	}

	return nil
}

func ComparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	return err == nil
}
