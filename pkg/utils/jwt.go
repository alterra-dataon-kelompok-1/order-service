package utils

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Password string    `json:"password,omitempty"`
	Role     Role      `json:"role"`
}

type Role string

const (
	Customer Role = "customer"
	Staff    Role = "staff"
	Admin    Role = "admin"
)

type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Role  Role   `json:"role"`
	jwt.StandardClaims
}

func GenerateJWTToken(user *User) (string, error) {
	claims := Claims{
		Email: user.Email,
		ID:    user.ID.String(),
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ExtractBearerToken(r *http.Request) string {
	headerAuthorization := r.Header.Get("Authorization")
	bearerToken := strings.Split(headerAuthorization, " ")
	return html.EscapeString(bearerToken[1])
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	jwtKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	return token, err
}

func ExtractJWTClaimsFromString(tokenString string) (*Claims, error) {
	jwtKey := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		fmt.Printf("%v %v %v", claims.Email, claims.Role, claims.StandardClaims.ExpiresAt)
	} else {
		fmt.Println(err)
	}

	return claims, nil
}
