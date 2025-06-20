package manager

import (
	"fmt"
	"strings"
)

// UrlManager gestisce la generazione degli URL per le immagini
type UrlManager struct {
	baseUrl string
}

// NewUrlManager crea una nuova istanza del manager URL
func NewUrlManager(baseUrl string) *UrlManager {
	// Rimuove il trailing slash se presente
	baseUrl = strings.TrimSuffix(baseUrl, "/")

	return &UrlManager{
		baseUrl: baseUrl,
	}
}

// GetImageUrl restituisce l'URL completo per un'immagine dato il nome del file
func (um *UrlManager) GetImageUrl(imageName string) string {
	return fmt.Sprintf("%s/media/%s", um.baseUrl, imageName)
}

// GetThumbnailUrl restituisce l'URL completo per un thumbnail dato il nome del file
func (um *UrlManager) GetThumbnailUrl(imageName string) string {
	return fmt.Sprintf("%s/media/thumbnails/%s", um.baseUrl, imageName)
}

// GetPreviewUrl restituisce l'URL completo per un'anteprima dato il nome del file
func (um *UrlManager) GetPreviewUrl(imageName string) string {
	return fmt.Sprintf("%s/media/previews/%s", um.baseUrl, imageName)
}
