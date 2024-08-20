package middleware

import (
	"context"
	"net/http"
	"strings"

	"Student_RESTAPI/utils"

	"github.com/dgrijalva/jwt-go"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		http.Error(w, "Unauthorized: Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := utils.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Unauthorized: Invalid token claims", http.StatusUnauthorized)
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Unauthorized: Invalid user ID in token", http.StatusUnauthorized)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "userID", userID))

	next(w, r)
}
