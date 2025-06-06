package controller

import (
	"net/http"

	"wedding-photo-backend/internal/weddingphoto/model"
	"wedding-photo-backend/internal/weddingphoto/service"

	"github.com/gin-gonic/gin"
)

// PhotoController gestisce le operazioni sulle foto
type PhotoController struct {
	photoService *service.PhotoService
}

// NewPhotoController crea una nuova istanza del controller
func NewPhotoController(photoService *service.PhotoService) *PhotoController {

	return &PhotoController{
		photoService: photoService,
	}
}

// UploadPhoto gestisce l'upload di una foto
// @Summary Upload di una foto
// @Description Carica una nuova foto sul server
// @Tags photos
// @Accept json
// @Produce json
// @Param request body model.AddPhotoRequest true "Dati della foto da caricare"
// @Success 200 {object} model.AddPhotoResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /api/photos [post]
func (pc *PhotoController) AddPhoto(c *gin.Context) {
	var request model.AddPhotoRequest

	// Binding del JSON body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: "Dati richiesta non validi: " + err.Error(),
		})
		return
	}

	// Salva la foto tramite il service
	photo, err := pc.photoService.AddPhoto(request.ImageContent, request.ImageName)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.AddPhotoResponse{
		Photo: *photo,
	})
}

// GetPhotos restituisce la lista delle foto
// @Summary Recupera la lista delle foto
// @Description Ottiene tutte le foto caricate sul server
// @Tags photos
// @Produce json
// @Success 200 {array} model.Photo
// @Failure 500 {object} model.ErrorResponse
// @Router /api/photos [get]
func (pc *PhotoController) GetPhotos(c *gin.Context) {
	photos, err := pc.photoService.GetPhotoList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Errore nel recupero delle foto: " + err.Error(),
		})
		return
	}

	getPhotosResponse := model.GetPhotosResponse{
		Photos: photos,
	}

	c.JSON(http.StatusOK, getPhotosResponse)
}

// SetupRoutes configura tutte le route relative alle foto
func (pc *PhotoController) SetupRoutes(api *gin.RouterGroup) {
	photos := api.Group("/photos")
	{
		photos.POST("", pc.AddPhoto)
		photos.GET("", pc.GetPhotos)
	}
}
