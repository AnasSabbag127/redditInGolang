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

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var newPost models.InputPost
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&newPost); err != nil {
		http.Error(w, "Invalid Request Body ", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDb()
	if err != nil {
		http.Error(w, "DATABASE connection ERROR ", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// insert new post
	_, err = db.Exec("INSERT INTO posts(post_title,post_text,user_id) VALUES($1,$2,$3) ", newPost.PostTitle, newPost.PostText, newPost.UserId)

	if err != nil {
		log.Println("error is : ", err)
		http.Error(w, "DATABASE ERROR ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "post created successful"}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Marshell error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	postId, err := uuid.Parse(idParam)

	if err != nil {
		http.Error(w, "invalid post id type UUID ", http.StatusBadRequest)
		return
	}

	var post models.Post
	db, err := database.ConnectToDb()
	if err != nil {
		http.Error(w, "connection to db error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.QueryRow("SELECT id,post_title,post_text,user_id FROM posts WHERE id = $1", postId).Scan(&post.Id, &post.PostTitle, &post.PostText, &post.UserId)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "post successfully created",
		"post":    post,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Marshal Error", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)

}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	IdParam := r.URL.Query().Get("id")
	postId, err := uuid.Parse(IdParam)
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

	_, err = db.Exec("DELETE FROM posts WHERE id = $1", postId)

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

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	postId, err := uuid.Parse(idParam)
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
	var post models.Post
	defer db.Close()
	err = db.QueryRow("SELECT user_id,post_title,post_text FROM posts WHERE id=$1", postId).Scan(&post.UserId, &post.PostTitle, &post.PostText)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Row Found ")
			http.Error(w, "post Not Found ", http.StatusNotFound)
			return
		}
		log.Println("Error is : ", err)
		http.Error(w, "DB Internal Error ", http.StatusInternalServerError)
		return
	}
	//now here its valid uuid for update post..
	var updatePost models.InputPost
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatePost); err != nil {
		log.Println("Error is : ", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	//now update the post query
	_, err = db.Exec("UPDATE posts SET post_title=$1,post_text=$2 WHERE id = $3", updatePost.PostTitle, updatePost.PostText, postId)
	if err != nil {
		log.Println("ERROR is : ", err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "update post success"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Marshal Error : ", err)
		http.Error(w, "Marshal Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}
