# User Management

使用者帳號管理系統，前端介面由 Go 撰寫，透過 REST API 呼叫 `user-data-service` 微服務取得與管理 users 資料。

## 📁 專案結構

user_management/
├── main.go                  # 啟動程式，啟用路由
├── handlers/
│   ├── auth.go              # 登入、註冊、驗證Token
│   ├── user_handler.go      # 使用者管理：列出所有使用者、單一使用者維護頁面、更新密碼
├── templates/
│   ├── login.html           # 登入頁面
│   ├── register.html        # 註冊頁面
│   ├── users.html           # 管理員用，列出所有使用者
│   ├── user_edit.html       # 維護使用者資料頁面(自己或管理者)
└── utils/
    └── paseto_client.go     # 呼叫 PASETO Auth Service
