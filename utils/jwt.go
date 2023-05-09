package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("TopSecret")

func GenerateToken(username string) (string, error) {
	exp := time.Now().Add(10 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":        exp,
		"authorized": true,
		"user":       "username",
	})

	fmt.Printf("Token: %v\n", token)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}

	fmt.Println("Token string", tokenString)
	return tokenString, nil
}

func ValidateToken(input string) (string, error) {
	validated, err := jwt.Parse(input, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if claims, ok := validated.Claims.(jwt.MapClaims); ok && validated.Valid {
		fmt.Println(claims["exp"])
	} else {
		fmt.Println(err)
	}

	fmt.Println("Validated", validated)

	return input, err
}
