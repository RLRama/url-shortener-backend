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

	ctx.SetCookie(&iris.Cookie{
		Name:     "auth_token",
		Value:    token,
		MaxAge:   int(time.Hour * 24 * 7 / time.Second), // 1 week
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: iris.SameSiteLaxMode,
	})

	ctx.StatusCode(iris.StatusOK)
	err3 := ctx.JSON(iris.Map{
		"token":   token,
		"message": "login successful",
	})
	if err3 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err3)
		return
	}
}

func handleLogout(ctx iris.Context) {
	ctx.SetCookie(&iris.Cookie{
		Name:     "auth_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: iris.SameSiteLaxMode,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	ctx.StatusCode(iris.StatusOK)
	err := ctx.JSON(iris.Map{
		"message": "logout successful",
	})
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}
}

func handleUpdatePassword(ctx iris.Context) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	var req PasswordUpdateRequest
	if err2 := ctx.ReadJSON(&req); err2 != nil {
		ctx.StopWithError(iris.StatusBadRequest, err2)
		return
	}

	if !verifyPassword(req.CurrentPassword, user.Password) {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		_, err2 := ctx.WriteString("invalid current password")
		if err2 != nil {
			return
		}
		return
	}

	if req.CurrentPassword == req.NewPassword {
		ctx.StopWithStatus(iris.StatusBadRequest)
		_, err2 := ctx.WriteString("new password cannot be same as current password")
		if err2 != nil {
			return
		}
		return
	}

	// new password validation
	if err2 := validateFieldLength(req.NewPassword, "new_password", 8, 100); err2 != nil {
		ctx.StopWithError(iris.StatusBadRequest, err2)
		return
	}
	if err2 := validatePasswordCharTypes(req.NewPassword); err2 != nil {
		ctx.StopWithError(iris.StatusBadRequest, err2)
		return
	}

	hashedPassword, err2 := hashPassword(req.NewPassword)
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	now := time.Now().UTC()
	userKey := fmt.Sprintf("user:%d", user.ID)
	_, err2 = rdb.HSet(ctx, userKey,
		"password", hashedPassword,
		"updated_at", now.Format(time.RFC3339),
	).Result()
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	_, err = ctx.WriteString("password updated")
	if err != nil {
		return
	}
}

func handleUpdateUsername(ctx iris.Context) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err)
		return
	}

	var req UsernameUpdateRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StopWithError(iris.StatusBadRequest, err)
		return
	}

	if !verifyPassword(req.Password, user.Password) {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		_, err2 := ctx.WriteString("invalid current password")
		if err2 != nil {
			return
		}
		return
	}

	if err2 := validateFieldLength(req.NewUsername, "new_username", 3, 50); err2 != nil {
		ctx.StopWithError(iris.StatusBadRequest, err2)
		return
	}

	exists, err2 := checkUsernameExists(req.NewUsername)
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}
	if exists {
		ctx.StopWithStatus(iris.StatusConflict)
		_, err2 := ctx.WriteString("username already exists")
		if err2 != nil {
			return
		}
		return
	}

	now := time.Now().UTC()
	userKey := fmt.Sprintf("user:%d", user.ID)
	_, err2 = rdb.HSet(ctx, userKey,
		"username", req.NewUsername,
		"updated_at", now.Format(time.RFC3339)).Result()
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	newToken, err2 := generateToken(req.NewUsername)
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	err3 := ctx.JSON(iris.Map{
		"message": "username updated",
		"token":   newToken,
	})
	if err3 != nil {
		return
	}
}
