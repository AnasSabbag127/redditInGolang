package api

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	var newUser models.InputUser

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&newUser); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	db, err := database.ConnectToDb()

	if err != nil {
		http.Error(w, "DB Err", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// insert new user
	_, err = db.Exec("INSERT INTO users(name,password) VALUES($1,$2)", newUser.Name, newUser.Password)

	if err != nil {
		http.Error(w, "Database Error ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	response := map[string]interface{}{"message": "User created successfully",
		"user": newUser,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "JSON marshaling error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

//for get user details ...

func GetUserHandler(w http.ResponseWriter, r *http.Request) {

	idParam := r.URL.Query().Get("id")
	userId, err := uuid.Parse(idParam) // convert into uuid

	if err != nil {
		log.Println("error is : ", err)
		http.Error(w, "Invalid UUID ", http.StatusBadRequest)
		return
	}
	log.Println("the uuid is ", userId)
	var user models.User
	db, err := database.ConnectToDb()

	if err != nil {
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
	}
	defer db.Close()

	err = db.QueryRow("SELECT id,name,password FROM users WHERE id = $1", userId).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found ", http.StatusNotFound)
			return
		}
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}
	// convert user struct to json and write to response
	response := map[string]interface{}{
		"message": "User retrieved successfully",
		"user":    user, // your retrieved user struct
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "JSON marshaling error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	IdParam := r.URL.Query().Get("id")
	userId, err := uuid.Parse(IdParam)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, "Invalid user uuid ", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDb()
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = $1", userId)

	if err != nil {
		log.Println("Error : ", err.Error())
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}
	response := map[string]string{"status": "success"}
	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.Println("Error : ", err.Error())
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	userId, err := uuid.Parse(idParam)
	if err != nil {
		log.Println("Error : ", err.Error())
		http.Error(w, "Invalid UUid", http.StatusInternalServerError)
		return
	}
	db, err := database.ConnectToDb()
	if err != nil {
		log.Println("Error : ", err.Error())
		http.Error(w, "DATABASE Connection ERROR", http.StatusInternalServerError)
		return
	}
	var user models.User
	defer db.Close()
	err = db.QueryRow("SELECT id,name,password FROM users WHERE id=$1", userId).Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Row Found ")
			http.Error(w, "user Not Found ", http.StatusNotFound)
			return
		}
		log.Println("Error is : ", err)
		http.Error(w, "DB Internal Error ", http.StatusInternalServerError)
		return
	}
	//now here its valid uuid for update user..
	var updateUser models.InputUser
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateUser); err != nil {
		log.Println("Error is : ", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	//now update the user query
	_, err = db.Exec("UPDATE users SET name=$1,password=$2 WHERE id = $3", updateUser.Name, updateUser.Password, userId)
	if err != nil {
		log.Println("ERROR is : ", err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "update success"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Marshal Error : ", err)
		http.Error(w, "Marshal Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}
