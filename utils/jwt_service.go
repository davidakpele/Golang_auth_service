package utils

import (
	"time"
	"os"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(userID uint, email string) (string, error) {
    claims := jwt.MapClaims{
        "id":    userID,
        "email": email,
        "roles": []string{"USER"},               // Example role claim
        "exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    return signedToken, nil
}