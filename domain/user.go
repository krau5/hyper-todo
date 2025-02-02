package domain

import "time"

type User struct {
	ID        int64     `json:"-" gorm:"unique;autoIncrement"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
