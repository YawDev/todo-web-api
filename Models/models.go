package models

import (
	"time"
)

type Task struct {
	Id          int       `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `gorm:"default:false" json:"IsCompleted"`
	ListId      int       `gorm:"foreignkey:ListId" json:"list_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type List struct {
	Id        int       `gorm:"primaryKey" json:"id"`
	Tasks     []Task    `json:"tasks"`
	UserId    int       `gorm:"foreignkey:UserId" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type User struct {
	Id        int       `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:100;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
