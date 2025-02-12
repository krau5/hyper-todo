package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/config"
	_ "github.com/krau5/hyper-todo/docs"
	"github.com/krau5/hyper-todo/internal/repository"
	"github.com/krau5/hyper-todo/internal/rest"
	"github.com/krau5/hyper-todo/task"
	"github.com/krau5/hyper-todo/user"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := config.GetDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to db:", err)
	}

	err = db.AutoMigrate(&repository.UserModel{}, &repository.TaskModel{})
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	log.Println("Migrations ran successfully")

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	usersRepo := repository.NewUserRepository(db)
	usersService := user.NewService(usersRepo)

	tasksRepo := repository.NewTasksRepository(db)
	tasksService := task.NewService(tasksRepo, usersRepo)

	rest.NewPingHandler(r)
	rest.NewAuthHandler(r, usersService)
	rest.NewTasksHandler(r, tasksService)
	rest.NewUsersHandler(r, usersService)

	r.Run()
}
