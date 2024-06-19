package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	"strings"
)

func authMiddleware(ctx iris.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			ctx.StopWithError(iris.StatusUnauthorized, err)
			return
		} else {
			ctx.StopWithError(iris.StatusBadRequest, err)
			return
		}
	}

	if !token.Valid {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		_, err2 := ctx.WriteString("invalid token")
		if err2 != nil {
			return
		}
		return
	}

	// if token is valid, refresh it
	newToken, err2 := generateToken(claims.Username)
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	ctx.Header("New-Token", newToken)
	ctx.Next()
}
