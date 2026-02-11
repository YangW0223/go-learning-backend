package model

import "time"

// User 是用户领域模型。
// 注意：PasswordHash 只用于持久化与校验，不应对外输出。
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}
