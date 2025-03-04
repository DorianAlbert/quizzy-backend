package quizzy

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"quizzy.app/backend/quizzy/cfg"
	quizzyhttp "quizzy.app/backend/quizzy/http"
	"quizzy.app/backend/quizzy/services"
	"time"
)

func Run() {
	config := cfg.LoadCfgFromEnv()

	// Configure GIN execution mode (dev, test, production).
	setGinMode(config.Env)

	log.Printf("application running in %s mode.\n", config.Env)

	// Initializing GIN engine.
	engine := gin.Default()

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Apply CORS middleware
	engine.Use(cors.New(corsConfig))

	router := engine.Group(config.BasePath)

	// Firebase and Redis middleware remain the same
	router.Use(func(ctx *gin.Context) {
		if client, err := services.ConfigureFirebase(config); err == nil {
			ctx.Set("firebase-services", client)
		}
	})
	router.Use(func(ctx *gin.Context) {
		redis := services.ConfigureRedis(config)
		if redis != nil {
			ctx.Set("redis-service", redis)
		}
	})

	// Initializing HTTP routes.
	quizzyhttp.ConfigureRouting(router)

	// Running server...
	if err := engine.Run(config.Addr); err != nil {
		log.Fatalf("Failed to start server on %s: %s", config.Addr, err)
	}
}

func setGinMode(env string) {
	switch env {
	case cfg.EnvDevelopment:
		gin.SetMode(gin.DebugMode)
	case cfg.EnvTest:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}
