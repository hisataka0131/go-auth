package model

import "github.com/jinzhu/gorm"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	gorm.Model
}

type JWT struct {
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}
