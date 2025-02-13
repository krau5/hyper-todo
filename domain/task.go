package domain

import "time"

type Task struct {
	ID          int64     `json:"id" gorm:"unique;autoIncrement"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Deadline    time.Time `json:"deadline"`
	Completed   bool      `json:"completed,omitempty" gorm:"default:false"`
	UserId      int64     `json:"-" gorm:"not null"`
}

type UpdateTaskData struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Completed   *bool      `json:"completed,omitempty"`
}
