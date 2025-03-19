package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func AuthenticationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Authorization header is missing"})
            c.Abort()
            return
        }

        // Extract token from "Bearer <token>" format
        parts := strings.Fields(authHeader)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token format"})
            c.Abort()
            return
        }
        tokenString := parts[1]

        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Ensure signing method matches
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            log.Printf("Token validation error: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Extract claims and set them in context
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token claims"})
            c.Abort()
            return
        }

        c.Set("user_id", claims["id"])
        c.Set("email", claims["email"])
        c.Set("roles", claims["roles"])

        c.Next()
    }
}