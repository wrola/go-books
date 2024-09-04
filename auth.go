package main

import ( 
	"net/http" 
 	"fmt" 
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func Authorize(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
    	return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Println("TOKEN", token)
}