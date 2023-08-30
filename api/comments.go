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

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var newComment models.InputComment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newComment); err != nil {
		log.Println("Error : ", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDb()
	if err != nil {
		log.Println("Error ", err)
		http.Error(w, "Database connection error ", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	//insert new comment
	_, err = db.Exec("INSERT INTO comments(user_id,post_id,comment) VALUES($1,$2,$3)", newComment.UserId, newComment.PostId, newComment.Comment)
	if err != nil {
		log.Println("Error is ", err)
		http.Error(w, "DATABASE ERROR ", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "success"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Marshal error ", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	commentId, err := uuid.Parse(idParam)
	if err != nil {
		log.Println("invalid path uuid ", err)
		http.Error(w, "Invalid uuid", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDb()
	if err != nil {
		log.Println("Error is ", err)
		http.Error(w, "DB CONNECTION ERROR ", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	var getComm models.GetComment
	//here get comment
	err = db.QueryRow("SELECT user_id,comment FROM comments WHERE id = $1", commentId).Scan(&getComm.UserId, &getComm.Comment)

	if err != nil {
		log.Println("Error is ", err)
		http.Error(w, "not found any row ", http.StatusNotFound)
		return
	}
	response := map[string]interface{}{
		"status":  "success",
		"comment": getComm,
	}

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		log.Println("Error is ", err)
		http.Error(w, "Marshal Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	IdParam := r.URL.Query().Get("id")
	commentId, err := uuid.Parse(IdParam)
	if err != nil {
		log.Println("Error: ", err)
		http.Error(w, "Invalid comment uuid ", http.StatusBadRequest)
		return
	}
	db, err := database.ConnectToDb()
	if err != nil {
		log.Println("Error:", err)
		http.Error(w, "failed to connect DB ", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM comments WHERE id = $1", commentId)

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

func UpdateCommentHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	commentId, err := uuid.Parse(idParam)
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
	var comment models.Comments
	defer db.Close()
	err = db.QueryRow("SELECT comment FROM comments WHERE id=$1", commentId).Scan(&comment.Comment)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No Row Found ")
			http.Error(w, "Comment Not Found ", http.StatusNotFound)
			return
		}
		log.Println("Error is : ", err)
		http.Error(w, "DB Internal Error ", http.StatusInternalServerError)
		return
	}
	//now here its valid uuid for update comment..
	var updateComment models.UpdateComment
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateComment); err != nil {
		log.Println("Error is : ", err)
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	//now update the comment query
	_, err = db.Exec("UPDATE comments SET comment=$1 WHERE id = $2", updateComment.Comment, commentId)
	if err != nil {
		log.Println("ERROR is : ", err)
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "update comment success"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Println("Marshal Error : ", err)
		http.Error(w, "Marshal Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)

}
