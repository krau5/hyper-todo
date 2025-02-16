package main

import (
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/krau5/hyper-todo/config"
	_ "github.com/krau5/hyper-todo/docs"
	"github.com/krau5/hyper-todo/internal/repository"
	"github.com/krau5/hyper-todo/internal/rest"
	"github.com/krau5/hyper-todo/task"
	"github.com/krau5/hyper-todo/user"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := c.Writer.Status()

		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), string(status)).Inc()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath(), string(status)).Observe(duration)
	}
}

// @title Hyper Todo API
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := initLogger(gin.Mode())
	defer logger.Sync()

	db := initDB(logger)

	r := gin.Default()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))
	registerHandlers(r, db)

	logger.Info("Server started", zap.String("port", config.Envs.Port))

	if err := r.Run(":" + config.Envs.Port); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

func initDB(logger *zap.Logger) *gorm.DB {
	dsn := config.GetDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to db", zap.Error(err))
	}

	err = db.AutoMigrate(&repository.UserModel{}, &repository.TaskModel{})
	if err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	return db
}

func initLogger(mode string) *zap.Logger {
	if mode == "release" {
		return zap.Must(zap.NewProduction())
	}
	return zap.Must(zap.NewDevelopment())
}

func registerHandlers(r *gin.Engine, db *gorm.DB) {
	usersRepo := repository.NewUserRepository(db)
	usersService := user.NewService(usersRepo)

	tasksRepo := repository.NewTasksRepository(db)
	tasksService := task.NewService(tasksRepo, usersRepo)

	r.Use(prometheusMiddleware())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	rest.NewPingHandler(r)
	rest.NewAuthHandler(r, usersService)
	rest.NewTasksHandler(r, tasksService)
	rest.NewUsersHandler(r, usersService)

	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.DefaultModelsExpandDepth(-1)))
}
