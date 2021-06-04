package model

import "time"

type USER struct {
	ID        string    `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	Username  string    `json:"username" sql:"username"`
	Role      string    `json:"role"`
	Level     string    `json:"level"`
	TokenHash string    `json:"tokenHash" sql:"tokenhash"`
	CreatedAt time.Time `json:"createDate" sql:"createAt"`
	UpdatedAt time.Time `json:"updateDate" sql:"updateAt"`
	LastLogin time.Time `json:"loginDate" sql:"loginAt"`
}

// USERS array of USER type
type USERS []USER
