package service

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"wedding-photo-backend/internal/weddingphoto/manager"
	"wedding-photo-backend/internal/weddingphoto/model"
)

// PhotoService gestisce la logica di business per le foto
type PhotoService struct {
	photoManager *manager.PhotoManager
	urlManager   *manager.UrlManager
	queueManager *manager.QueueManager
}

// NewPhotoService crea una nuova istanza del service
func NewPhotoService(photoManager *manager.PhotoManager, urlManager *manager.UrlManager, queueManager *manager.QueueManager) *PhotoService {
	return &PhotoService{
		photoManager: photoManager,
		urlManager:   urlManager,
		queueManager: queueManager,
	}
}

// GetPhotoList restituisce la lista delle immagini salvate con paginazione
func (ps *PhotoService) GetPhotoList(page, perPage int) ([]model.Photo, int, error) {
	// Recupera la lista delle immagini dal manager
	imageNames, err := ps.photoManager.GetPhotoList()
	if err != nil {
		return nil, 0, fmt.Errorf("errore nel recupero della lista delle immagini: %v", err)
	}

	var photos []model.Photo
	for _, imageName := range imageNames {
		// Verifica se thumbnail e preview esistono
		if !ps.photoManager.ThumbnailExists(imageName) || !ps.photoManager.PreviewExists(imageName) {
			// Salta la foto se thumbnail o preview non esistono
			continue
		}

		thumbnailUrl := ps.urlManager.GetThumbnailUrl(imageName)
		previewUrl := ps.urlManager.GetPreviewUrl(imageName)

		photos = append(photos, model.Photo{
			ImageName: imageName,
			// ImageUrl:     ps.urlManager.GetImageUrl(imageName),
			ThumbnailUrl: thumbnailUrl,
			PreviewUrl:   previewUrl,
		})
	}

	// Ordina le foto per nome file in ordine decrescente (assumendo che i nomi file contengano timestamp)
	sort.Slice(photos, func(i, j int) bool {
		return photos[i].ImageName > photos[j].ImageName
	})

	// Calcola il numero totale di pagine
	totalPhotos := len(photos)
	totalPages := int(math.Ceil(float64(totalPhotos) / float64(perPage)))

	// Calcola gli indici per la paginazione
	startIndex := (page - 1) * perPage
	endIndex := startIndex + perPage

	// Verifica i limiti
	if startIndex >= totalPhotos {
		return []model.Photo{}, totalPages, nil
	}

	if endIndex > totalPhotos {
		endIndex = totalPhotos
	}

	// Ritorna la finestra di foto richiesta
	paginatedPhotos := photos[startIndex:endIndex]

	return paginatedPhotos, totalPages, nil
}

// AddPhoto salva una foto da multipart form data e aggiunge alla coda di elaborazione
func (ps *PhotoService) AddPhoto(fileReader io.Reader, imageName string, contentType string, fileSize int64) (*model.Photo, error) {
	// Rileva il MIME type reale dal contenuto del file
	realMimeType, newReader, err := ps.photoManager.DetectMimeTypeFromBytes(fileReader)
	if err != nil {
		return nil, fmt.Errorf("errore nella lettura del file: %v", err)
	}

	// Verifica che il MIME type reale sia un'immagine supportata
	if !ps.photoManager.IsValidImageMimeType(realMimeType) {
		return nil, fmt.Errorf("il file non è un'immagine valida o il formato non è supportato")
	}

	// Verifica che il MIME type dichiarato corrisponda a quello reale (opzionale, per maggiore sicurezza)
	if !ps.isImageMimeType(contentType) || !ps.mimeTypesMatch(contentType, realMimeType) {
		// Log di warning ma usa il MIME type reale
		fmt.Printf("Warning: MIME type dichiarato (%s) diverso da quello reale (%s)\n", contentType, realMimeType)
	}

	// Usa il MIME type reale per il salvataggio
	fileName, err := ps.photoManager.SavePhotoFromBytes(newReader, imageName, realMimeType, fileSize)
	if err != nil {
		return nil, err
	}

	// Aggiunge l'immagine alla coda di elaborazione
	if err := ps.queueManager.AddImageToQueue(fileName); err != nil {
		fmt.Printf("Errore nell'aggiunta dell'immagine alla coda: %v\n", err)
		// Non restituiamo errore, il file è stato comunque salvato
	}

	// Crea e restituisce l'oggetto Photo con URL completo
	photo := &model.Photo{
		ImageName:    fileName,
		ImageUrl:     ps.urlManager.GetImageUrl(fileName),
		ThumbnailUrl: ps.urlManager.GetThumbnailUrl(fileName),
		PreviewUrl:   ps.urlManager.GetPreviewUrl(fileName),
	}

	return photo, nil
}

// mimeTypesMatch verifica se i MIME types sono compatibili
func (ps *PhotoService) mimeTypesMatch(declared, real string) bool {
	// Normalizza i MIME types
	if declared == "image/jpg" {
		declared = "image/jpeg"
	}

	return declared == real || (declared == "image/jpeg" && real == "image/jpeg")
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
