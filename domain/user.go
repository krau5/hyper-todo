package domain

type User struct {
	ID       int64  `json:"-" gorm:"unique;autoIncrement"`
	Name     string `json:"name" gorm:"not null" example:"user"`
	Email    string `json:"email" gorm:"unique;not null" example:"user@example.com"`
	Password string `json:"-" gorm:"not null"`
}
