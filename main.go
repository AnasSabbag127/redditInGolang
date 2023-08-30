package main

import (
	"awesomeProject/api"
	"awesomeProject/database"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func handler(_ http.ResponseWriter, r *http.Request) {
	log.Println("URL Path=%q ", r.URL.Path)
}
func main() {

	fmt.Println("hello World!")
	_, err := database.ConnectToDb()
	if err != nil {
		fmt.Println("db connection error")
		return
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/create_users", api.CreateUserHandler)
	http.HandleFunc("/api/get_users", api.GetUserHandler)
	http.HandleFunc("/api/update_user", api.UpdateUserHandler)
	http.HandleFunc("/api/delete_user", api.DeleteUserHandler)

	http.HandleFunc("/api/create_posts", api.CreatePostHandler)
	http.HandleFunc("/api/get_posts", api.GetPostHandler)
	http.HandleFunc("/api/update_post", api.UpdatePostHandler)
	http.HandleFunc("/api/delete_post", api.DeletePostHandler)

	http.HandleFunc("/api/create_comments", api.CreateCommentHandler)
	http.HandleFunc("/api/get_comments", api.GetCommentHandler)
	http.HandleFunc("/api/update_comment", api.UpdateCommentHandler)
	http.HandleFunc("/api/delete_comment", api.DeleteCommentHandler)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))

}
