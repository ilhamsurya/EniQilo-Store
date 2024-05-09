package entity

import (
	"database/sql"
	"time"
)

type User struct {
	UserId      uint32       `db:"user_id" json:"userId"`
	Email       string       `db:"email" json:"email"`
	Name        string       `db:"name" json:"name"`
	PhoneNumber string       `db:"phone_number" json:"phoneNumber"`
	Password    string       `db:"password" json:"-"`
	Salt        string       `db:"salt" json:"-"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at" json:"updated_at"`
}

type UserParam struct {
	Email       string
	Name        string
	PhoneNumber string
	Password    string
	Salt        string
}

type UserLoginParam struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type UserResponse struct {
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	AccessToken string `json:"accessToken"`
}
