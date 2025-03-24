package models

import "time"

type User struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	LoginID   string    `json:"login_id"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	RoleID    int       `json:"role_id"`
	RoleName  *string   `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}
