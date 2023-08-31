package authUser

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("secret-key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Println("Error is : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userAcc models.User
	db, err := database.ConnectToDb()
	if err != nil {
		http.Error(w, "not connect to database ", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	err = db.QueryRow("SELECT password FROM users where name = $1", credentials.Username).Scan(&userAcc.Password)

	if err != nil {
		//use log info for error
		log.Println("Error is :", err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	// here change
	hashBytePassword := []byte(userAcc.Password)
	err = bcrypt.CompareHashAndPassword(hashBytePassword, []byte(credentials.Password))
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, "Invalid credentials.. ", http.StatusUnauthorized)
		return
	}

	// end ..
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}

func CheckForTokenValidation(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("Error is : ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	tokenStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil || !tkn.Valid {
		if err == jwt.ErrSignatureInvalid || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Println("Error is : ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenStr := cookie.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil || !tkn.Valid {
		if err == jwt.ErrSignatureInvalid || !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprintf("hello ,%s", claims.Username)))
}
