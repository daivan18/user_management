package models

import "time"

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"` // 原始密碼（來自註冊/登入）
	PasswordHash  string    `json:"-"`        // 雜湊後的密碼（寫入 DB）
	CellPhone     string    `json:"cell_phone"`
	CellPhoneHash string    `json:"cell_phone_hash"` // 新增
	LineID        string    `json:"line_id"`
	CreateTime    time.Time `json:"create_time"`
	UpdateTime    time.Time `json:"update_time"`
}
