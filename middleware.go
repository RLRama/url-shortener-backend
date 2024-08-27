package main

import (
	"errors"
	"fmt"
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

	user, err2 := getUserByUsername(claims.Username)
	if err2 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err2)
		return
	}

	issuedAt, err3 := claims.GetIssuedAt()
	if err3 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, fmt.Errorf("error getting token issue time: %w", err3))
		return
	}

	if issuedAt == nil || issuedAt.Before(user.UpdatedAt) {
		ctx.StopWithStatus(iris.StatusUnauthorized)
		_, err4 := ctx.WriteString("token invalidated due to account update")
		if err4 != nil {
			return
		}
		return
	}

	// store username in the context
	ctx.Values().Set("username", claims.Username)

	// if token is valid, refresh it
	newToken, err4 := generateToken(claims.Username)
	if err4 != nil {
		ctx.StopWithError(iris.StatusInternalServerError, err4)
		return
	}

	ctx.Header("New-Token", newToken)
	ctx.Next()
}
