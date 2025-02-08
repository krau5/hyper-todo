package domain

type User struct {
	ID       int64  `json:"-" gorm:"unique;autoIncrement"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
}
