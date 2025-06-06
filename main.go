package main

import (
	"log"
	_ "wedding-photo-backend/docs"
	"wedding-photo-backend/internal/weddingphoto/controller"
	"wedding-photo-backend/internal/weddingphoto/manager"
	"wedding-photo-backend/internal/weddingphoto/service"
	"wedding-photo-backend/internal/weddingphoto/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	_ = godotenv.Load()

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

	host := util.GetEnv("HOST", "0.0.0.0")
	port := util.GetEnv("PORT", "8739")
	baseUrl := util.GetEnv("BASE_URL", "http://localhost:8739")
	photosDir := util.GetEnv("PHOTOS_DIR", "media")

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

	// Avvia il server sulla porta
	log.Println("Server avviato su http://" + host + ":" + port)
	if err := r.Run(host + ":" + port); err != nil {
		log.Fatal("Errore nell'avvio del server:", err)
	}
}
