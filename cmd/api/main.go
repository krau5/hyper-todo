package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/config"
	"github.com/krau5/hyper-todo/internal/repository"
	"github.com/krau5/hyper-todo/internal/rest"
	"github.com/krau5/hyper-todo/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := config.GetDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	err = db.AutoMigrate(&repository.UserModel{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	} else {
		log.Println("Migrations ran successfully")
	}

	r := gin.Default()

	usersRepo := repository.NewUserRepository(db)
	usersService := user.NewService(usersRepo)

	rest.NewPingHandler(r)
	rest.NewUserHandler(r, usersService)

	r.Run()
}
