package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
