package main

import (
	"crypto/rand"
	"encoding/base64"
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

	dbUser, err := FindUserByUsername(authUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if dbUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": dbUser.Username})
}

func UpdateUsernameHandler(c *gin.Context) {
	user, _ := c.Get("user")
	authUser := user.(*User)

	userID, err := FindUserIDByUsername(authUser.Username)

	fmt.Println(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	var updatedUsername struct {
		Username string `json:"username"`
	}

	if err2 := c.ShouldBind(&updatedUsername); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	if updatedUsername.Username != "" && updatedUsername.Username != authUser.Username {
		exists, err2 := UsernameExists(updatedUsername.Username)
		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
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

	if err2 := UpdateUsernameInRedis(userID, updatedUsername.Username); err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
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
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("generated token: ", tokenString)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func UpdatePasswordHandler(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, _ := c.Get("user")
	user := authUser.(*User)
	userID, err := FindUserIDByUsername(user.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if dbUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
		return
	}

	if err2 := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.CurrentPassword)); err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
		return
	}

	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err2 := UpdatePasswordInRedis(userID, hashedPassword); err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
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
	salt, err := GenerateSalt()
	if err != nil {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func UpdateUsernameInRedis(userID string, newUsername string) error {
	userKey := fmt.Sprintf("user:%s", userID)

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

func UpdatePasswordInRedis(userID string, newPassword string) error {
	userKey := fmt.Sprintf("user:%s", userID)

	exists, err := rdb.Exists(ctx, userKey).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		return fmt.Errorf("user with ID %s does not exist", userID)
	}

	if err := rdb.HSet(ctx, userKey, "password", newPassword).Err(); err != nil {
		return err
	}

	return nil
}

func ComparePasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	return err == nil
}

func FindUserIDByUsername(username string) (string, error) {
	keys, err := rdb.Keys(ctx, "user:*").Result()
	if err != nil {
		return "", err
	}

	for _, key := range keys {
		user, err2 := rdb.HGetAll(ctx, key).Result()
		if err2 != nil {
			return "", err2
		}

		if user["username"] == username {
			userID := strings.TrimPrefix(key, "user:")
			return userID, nil
		}
	}

	return "", nil
}

func GetUserByID(userID string) (*User, error) {
	userKey := fmt.Sprintf("user:%s", userID)

	userData, err := rdb.HGetAll(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	if len(userData) == 0 {
		return nil, nil
	}

	user := &User{
		Username: userData["username"],
		Password: userData["password"],
	}

	return user, nil
}

func GenerateSalt() (string, error) {
	saltBytes := make([]byte, 16)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(saltBytes), nil
}
