// package main

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/dgrijalva/jwt-go"
// )

// func main() {
// 	http.HandleFunc("/validate-token", func(w http.ResponseWriter, r *http.Request) {
// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString == "" {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			fmt.Fprintln(w, "No token provided")
// 			return
// 		}

// 		// Parse and validate the token
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			// Check the signing method
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			// Provide the key used to sign the token (should match Django's settings)
// 			return []byte("your_secret_key"), nil
// 		})

// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			fmt.Fprintf(w, "Token parsing error: %v", err)
// 			return
// 		}

// 		if token.Valid {
// 			// Token is valid, you can process the request
// 			w.WriteHeader(http.StatusOK)
// 			fmt.Fprintln(w, "Token is valid")
// 		} else {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			fmt.Fprintln(w, "Invalid token")
// 		}
// 	})

// 	http.ListenAndServe(":8080", nil)
// }
