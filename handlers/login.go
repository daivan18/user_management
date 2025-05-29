// login.go
package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

// GET /login 顯示登入頁面
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))

	// 從 query 參數中讀取錯誤訊息
	errorMsg := r.URL.Query().Get("error")

	if cookie, err := r.Cookie("login_error"); err == nil {
		decoded, _ := url.QueryUnescape(cookie.Value)
		errorMsg = decoded

		// 刪除 cookie（設過期）
		http.SetCookie(w, &http.Cookie{
			Name:   "login_error",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}

	tmpl.Execute(w, map[string]string{
		"Error": errorMsg,
	})
}

// POST /login 處理登入邏輯
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// 失敗用 redirect 並帶錯誤訊息 query
		http.Redirect(w, r, "/login?error=表單解析失敗", http.StatusSeeOther)
		return
	}

	cellPhone := strings.TrimSpace(r.FormValue("cell_phone"))
	password := r.FormValue("password")

	if cellPhone == "" || password == "" {
		http.Redirect(w, r, "/login?error=請輸入手機號碼與密碼", http.StatusSeeOther)
		return
	}

	loginReq := map[string]string{
		"cell_phone": cellPhone,
		"password":   password,
	}
	bodyBytes, _ := json.Marshal(loginReq)

	resp, err := http.Post("http://localhost:8001/login", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil || resp.StatusCode != http.StatusOK {
		http.SetCookie(w, &http.Cookie{
			Name:  "login_error",
			Value: url.QueryEscape("登入失敗，請確認帳號密碼"),
			Path:  "/",
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		// http.Redirect(w, r, "/login?error=登入失敗，請確認帳號密碼", http.StatusSeeOther)
		return
	}
	defer resp.Body.Close()

	var loginResp struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		http.Redirect(w, r, "/login?error=登入回應解析失敗", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: loginResp.Token,
		Path:  "/",
	})

	// 登入成功也要 redirect
	if loginResp.Username == "admin" {
		http.Redirect(w, r, "/users", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/user/%s", loginResp.Username), http.StatusSeeOther)
	}
}

// 顯示登入頁面並呈現錯誤訊息
func renderLogin(w http.ResponseWriter, errMsg string) {
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, map[string]string{"Error": errMsg})
}
