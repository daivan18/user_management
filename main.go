package main

import (
	"log"
	"net/http"
	"os" // 引入 os 套件以讀取環境變數

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

	// 優先讀取 Cloud Run 分配的 PORT 環境變數，若無則預設 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server start at :%s\n", port)
	
	// 將原本寫死的 ":8080" 替換為動態的 ":" + port
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}