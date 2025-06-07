package controller

import (
	"net/http"
	"strconv"

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

// GetPhotos restituisce la lista delle foto con paginazione
// @Summary Recupera la lista delle foto
// @Description Ottiene tutte le foto caricate sul server con paginazione
// @Tags photos
// @Produce json
// @Param page query int false "Numero pagina (default: 1)"
// @Param perPage query int false "Elementi per pagina (default: 10, max: 100)"
// @Success 200 {object} model.GetPhotosResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/photos [get]
func (pc *PhotoController) GetPhotos(c *gin.Context) {
	// Parse dei parametri di paginazione
	page := 1
	perPage := 10

	if pageParam := c.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if perPageParam := c.Query("perPage"); perPageParam != "" {
		if pp, err := strconv.Atoi(perPageParam); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	photos, totalPages, err := pc.photoService.GetPhotoList(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Errore nel recupero delle foto: " + err.Error(),
		})
		return
	}

	getPhotosResponse := model.GetPhotosResponse{
		Photos:     photos,
		Page:       page,
		TotalPages: totalPages,
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
