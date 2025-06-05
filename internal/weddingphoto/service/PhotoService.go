package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"wedding-photo-backend/internal/weddingphoto/manager"
	"wedding-photo-backend/internal/weddingphoto/model"
)

// PhotoService gestisce la logica di business per le foto
type PhotoService struct {
	photoManager *manager.PhotoManager
	urlManager   *manager.UrlManager
}

// NewPhotoService crea una nuova istanza del service
func NewPhotoService(photoManager *manager.PhotoManager, urlManager *manager.UrlManager) *PhotoService {
	return &PhotoService{
		photoManager: photoManager,
		urlManager:   urlManager,
	}
}

// GetImageList restituisce la lista delle immagini salvate
func (ps *PhotoService) GetImageList() ([]model.Photo, error) {
	// Recupera la lista delle immagini dal manager
	imageNames, err := ps.photoManager.GetImageList()
	if err != nil {
		return nil, fmt.Errorf("errore nel recupero della lista delle immagini: %v", err)
	}

	var photos []model.Photo
	for _, imageName := range imageNames {
		photos = append(photos, model.Photo{
			ImageName: imageName,
			ImageUrl:  ps.urlManager.GetImageUrl(imageName),
		})
	}

	return photos, nil
}

// AddPhoto salva una foto da AddPhotoRequest e aggiorna la lista
func (ps *PhotoService) AddPhoto(imageContent string, imageName string) (*model.Photo, error) {
	// Decodifica il base64
	imageData, err := base64.StdEncoding.DecodeString(imageContent)
	if err != nil {
		return nil, fmt.Errorf("errore nella decodifica base64: %v", err)
	}

	// Determina il tipo MIME dall'header dei dati
	contentType := http.DetectContentType(imageData)
	if !ps.isImageMimeType(contentType) {
		return nil, fmt.Errorf("Il file deve essere un'immagine (JPEG, PNG, GIF, WebP)")
	}

	// Crea un reader dai dati dell'immagine DOPO la validazione
	imageReader := bytes.NewReader(imageData)

	// Salva tramite il manager
	fileName, err := ps.photoManager.SavePhotoFromBytes(imageReader, imageName, contentType, int64(len(imageData)))
	if err != nil {
		return nil, err
	}

	// Crea e restituisce l'oggetto Photo con URL completo
	photo := &model.Photo{
		ImageName: fileName,
		ImageUrl:  ps.urlManager.GetImageUrl(fileName),
	}

	return photo, nil
}

// isImageMimeType verifica se il MIME type è di un'immagine
func (ps *PhotoService) isImageMimeType(mimeType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if strings.HasPrefix(mimeType, validType) {
			return true
		}
	}
	return false
}
