package main

import (
	"log"
	"net/http"

	"user_management/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	r.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	r.HandleFunc("/register", handlers.RegisterPage).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")

	r.HandleFunc("/users", handlers.UsersPage).Methods("GET") // 管理員專用
	r.HandleFunc("/user/{username}", handlers.UserEditPage).Methods("GET")
	r.HandleFunc("/user/{username}", handlers.UserEditPost).Methods("POST")

	log.Println("Server start at :8080")
	http.ListenAndServe(":8080", r)
}
