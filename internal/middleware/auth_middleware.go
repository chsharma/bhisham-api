package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Secret key for verifying the JWT token
var jwtSecretKey = []byte("your_secret_key")

// JWTAuthentication middleware function
func JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Get JWT secret key from environment
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			panic("JWT_SECRET_KEY is not set in the environment")
		}
		jwtSecretKey := []byte(secretKey)

		// Token format: "Bearer <token>"
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := splitToken[1]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is correct
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecretKey, nil
		})

		// Handle token parsing errors
		if err != nil {
			if validationErr, ok := err.(*jwt.ValidationError); ok {
				if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
					// Token is expired
					http.Error(w, "Token expired", http.StatusForbidden)
					return
				} else if validationErr.Errors&(jwt.ValidationErrorMalformed|jwt.ValidationErrorUnverifiable|jwt.ValidationErrorSignatureInvalid) != 0 {
					// Invalid token
					http.Error(w, "Invalid token", http.StatusUnauthorized)
					return
				}
			}
			// Generic error for unhandled cases
			http.Error(w, "Error processing token", http.StatusInternalServerError)
			return
		}

		// If the token is not valid
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims if needed
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Example: Extract the user ID
			if userID, ok := claims["user_id"].(string); ok {
				r.Header.Set("User-ID", userID)
			}
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func GenerateNewToken(userID string) (string, error) {
	// Get JWT secret key from environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET is not set in the environment")
	}
	jwtSecretKey = []byte(secretKey)
	// Create a new JWT token with claims
	claims := jwt.MapClaims{
		"user_id": userID,                                     // User's ID
		"exp":     time.Now().Add(time.Hour * 24 * 90).Unix(), // Expiry time: 30 days
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
