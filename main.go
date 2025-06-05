package main

import (
	"log"
	_ "wedding-photo-backend/docs"
	"wedding-photo-backend/internal/weddingphoto/controller"
	"wedding-photo-backend/internal/weddingphoto/manager"
	"wedding-photo-backend/internal/weddingphoto/service"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Wedding Photo Backend API
// @version 1.0
// @description API per la gestione delle foto del matrimonio
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

func main() {
	// Inizializza il router Gin
	r := gin.Default()

	// Abilita CORS per consentire richieste da frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	photosDir := "media"
	baseUrl := "http://localhost:8080"

	photoManager := manager.NewPhotoManager(photosDir)
	urlManager := manager.NewUrlManager(baseUrl)
	photoService := service.NewPhotoService(photoManager, urlManager)
	photoController := controller.NewPhotoController(photoService)

	// Definisce le route API
	api := r.Group("/api")
	photoController.SetupRoutes(api)

	// Route per Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Set route /media as static file server
	r.Static("/media", photosDir)

	// Avvia il server sulla porta 8080
	log.Println("Server avviato su http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Errore nell'avvio del server:", err)
	}
}
