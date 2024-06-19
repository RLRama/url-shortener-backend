package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"time"
)

func helloWorldTest(ctx iris.Context) {
	_, err := ctx.WriteString("Hello World")
	if err != nil {
		ctx.StopWithStatus(iris.StatusInternalServerError)
		return
	}
}

func redisTest(ctx iris.Context) {
	err := rdb.Set(ctx, "key0", "dickson", 0).Err()
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	value, err4 := rdb.Get(ctx, "key0").Result()
	if err4 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err4)
		return
	}

	err6 := ctx.JSON(value)
	if err6 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err6)
		return
	}
}

func handleUserRegistration(ctx iris.Context) {
	var req RegisterUserRequest

	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	// fields validation
	if err := validateFieldLength(req.Username, "username", 3, 50); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	if err := validateFieldLength(req.Password, "password", 8, 100); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	if err := validatePasswordCharTypes(req.Password); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	// username existence check
	exists, err := checkUsernameExists(req.Username)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	if exists {
		ctx.StopWithStatus(iris.StatusConflict)
		return
	}

	hashedPassword, err2 := hashPassword(req.Password)
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	userID, err3 := rdb.Incr(ctx, "next_user_id").Result()
	if err3 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err3)
		return
	}

	now := time.Now().UTC()
	user := User{
		ID:        uint64(userID),
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// store user in Redis hash
	userKey := fmt.Sprintf("user:%d", user.ID)
	_, err4 := rdb.HSet(ctx, userKey, map[string]interface{}{
		"username":   user.Username,
		"password":   user.Password,
		"created_at": user.CreatedAt.Format(time.RFC3339),
		"updated_at": user.UpdatedAt.Format(time.RFC3339),
	}).Result()

	if err4 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err4)
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	err5 := ctx.JSON(iris.Map{
		"id":       user.ID,
		"username": user.Username,
	})
	if err5 != nil {
		return
	}
}

func handleLogin(ctx iris.Context) {
	var req RegisterUserRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	// find username and get his password
	user, err := getUserByUsername(req.Username)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
	if &user == nil {
		ctx.StopWithStatus(iris.StatusNotFound)
		_, err2 := ctx.WriteString("invalid username or password")
		if err2 != nil {
			return
		}
		return
	}

	if !verifyPassword(req.Password, user.Password) {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		_, err2 := ctx.WriteString("invalid username or password")
		if err2 != nil {
			return
		}
		return
	}

	token, err2 := generateToken(req.Username)
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	err3 := ctx.JSON(iris.Map{
		"token": token,
	})
	if err3 != nil {
		return
	}
}
