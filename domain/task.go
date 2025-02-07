package domain

import "time"

type Task struct {
	ID          int64     `json:"id" gorm:"unique;autoIncrement"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Deadline    time.Time `json:"deadline"`
	Completed   bool      `json:"completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	UserId      int64     `json:"-" gorm:"not null"`
}
