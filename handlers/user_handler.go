package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"user_management/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GET /user/{username}
func UserEditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	resp, err := http.Get(fmt.Sprintf("http://localhost:8001/users/%s", username))
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "無法取得使用者資料", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		http.Error(w, "回傳資料格式錯誤", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/user_edit.html"))
	tmpl.Execute(w, user)
}

// POST /user/{username}
func UserEditPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if err := r.ParseForm(); err != nil {
		http.Error(w, "表單解析失敗", http.StatusBadRequest)
		return
	}

	newPassword := r.FormValue("password")
	confirmPassword := r.FormValue("password2")

	if newPassword == "" || confirmPassword == "" || newPassword != confirmPassword {
		http.Error(w, "密碼欄位錯誤", http.StatusBadRequest)
		return
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "密碼加密失敗", http.StatusInternalServerError)
		return
	}

	payload := map[string]string{
		"password_hash": string(pwHash),
	}

	jsonBody, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8001/users/%s", username), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "密碼更新失敗", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/%s", username), http.StatusSeeOther)
}

// POST /user/{username}/delete
func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	req, _ := http.NewRequest("DELETE", "http://localhost:8001/users/"+username, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "刪除失敗", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// GET /users（admin only）
func UsersPage(w http.ResponseWriter, r *http.Request) {
	// TODO: 從 Redis 驗證 token，確認 admin 身分

	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	resp, err := http.Get("http://localhost:8001/users")
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "無法取得使用者資料", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var users []models.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		http.Error(w, "解析使用者資料失敗", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/users.html"))
	tmpl.Execute(w, users)
}

// GET /register
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/register.html"))
	tmpl.Execute(w, nil)
}

// POST /register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		showAlertAndRedirect(w, "表單解析失敗", "/register")
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	cellPhone := strings.TrimSpace(r.FormValue("cell_phone"))
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	if username == "" || cellPhone == "" || password == "" || password2 == "" {
		showAlertAndRedirect(w, "所有欄位皆需填寫", "/register")
		return
	}
	if password != password2 {
		showAlertAndRedirect(w, "兩次密碼不一致", "/register")
		return
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		showAlertAndRedirect(w, "密碼加密失敗", "/register")
		return
	}

	newUser := map[string]string{
		"username":      username,
		"password_hash": string(pwHash),
		"cell_phone":    cellPhone,
	}

	jsonBody, _ := json.Marshal(newUser)
	resp, err := http.Post("http://localhost:8001/users", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		showAlertAndRedirect(w, "註冊失敗，請稍後再試", "/register")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// 嘗試解析錯誤訊息（假設 user-data-service 有回傳 JSON 格式錯誤）
		var errResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errResp)
		msg := errResp["error"]
		if msg == "" {
			msg = "註冊失敗，帳號可能已存在"
		}
		showAlertAndRedirect(w, msg, "/register")
		return
	}

	showAlertAndRedirect(w, "註冊成功，請重新登入", "/login")
}

func showAlertAndRedirect(w http.ResponseWriter, message, redirectURL string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html>
		<head><meta charset="UTF-8"><title>通知</title></head>
		<body>
			<script>
				alert(%q);
				window.location.href = %q;
			</script>
		</body>
		</html>
	`, message, redirectURL)
}
