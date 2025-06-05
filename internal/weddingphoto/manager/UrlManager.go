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
