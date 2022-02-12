package entities

import "time"

type User struct {
	Id        int64     `form:"id" json:"id" query:"id"`
	Username  string    `form:"username" json:"username" query:"username"`
	Password  []byte    `form:"password" json:"password" query:"password"`
	UpdatedAt time.Time `form:"updated_at" json:"updated_at" query:"updated_at"`
	CreatedAt time.Time `form:"created_at" json:"created_at" query:"created_at"`
}
