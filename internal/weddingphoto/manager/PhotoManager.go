package manager

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// PhotoManager gestisce le operazioni sulfilesystem per le foto
type PhotoManager struct {
	photosDir     string
	thumbnailsDir string
}

// NewPhotoManager crea una nuova istanza del manager
func NewPhotoManager(photosDir string) *PhotoManager {
	thumbnailsDir := filepath.Join(photosDir, "thumbnails")

	// Crea le directory se non esistono
	if err := os.MkdirAll(photosDir, 0755); err != nil {
		fmt.Printf("Errore nella creazione della directory %s: %v\n", photosDir, err)
	}
	if err := os.MkdirAll(thumbnailsDir, 0755); err != nil {
		fmt.Printf("Errore nella creazione della directory thumbnails %s: %v\n", thumbnailsDir, err)
	}

	return &PhotoManager{
		photosDir:     photosDir,
		thumbnailsDir: thumbnailsDir,
	}
}

func (pm *PhotoManager) GetPhotoList() ([]string, error) {
	var images []string

	// Legge i file nella directory delle foto
	files, err := os.ReadDir(pm.photosDir)
	if err != nil {
		return nil, fmt.Errorf("errore nella lettura della directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && pm.isImageFile(file.Name()) {
			images = append(images, file.Name())
		}
	}

	return images, nil
}

// SavePhotoFromBytes salva una nuova immagine dal reader di bytes
func (pm *PhotoManager) SavePhotoFromBytes(reader io.Reader, originalFilename string, contentType string, size int64) (string, error) {
	// Genera un nome file unico con formato yyyy-mm-dd-hh-ii-ss-rand(0,99999999)
	now := time.Now()
	randomNum := rand.Intn(100000000) // 0-99999999

	// Estrae l'estensione dal nome originale o dal content type
	ext := filepath.Ext(originalFilename)
	if ext == "" {
		ext = pm.getExtensionFromMimeType(contentType)
	}

	filename := fmt.Sprintf("%04d-%02d-%02d-%02d-%02d-%02d-%08d%s",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		randomNum, ext)
	filePath := filepath.Join(pm.photosDir, filename)

	// Crea il file di destinazione
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("errore nella creazione del file: %v", err)
	}
	defer dst.Close()

	// Copia il contenuto dal reader al file
	_, err = io.Copy(dst, reader)
	if err != nil {
		return "", fmt.Errorf("errore nella scrittura del file: %v", err)
	}

	// Crea il thumbnail
	if err := pm.createThumbnail(filePath, filename, contentType); err != nil {
		fmt.Printf("Errore nella creazione del thumbnail per %s: %v\n", filename, err)
		// Non restituiamo errore per il thumbnail, continuiamo
	}

	return filename, nil
}

// createThumbnail crea un thumbnail di 400x400px usando ImageMagick
func (pm *PhotoManager) createThumbnail(originalPath, filename, contentType string) error {
	thumbnailPath := filepath.Join(pm.thumbnailsDir, filename)

	cmd := exec.Command("magick", originalPath,
		"-auto-orient",
		"-resize", "400x400^",
		"-gravity", "center",
		"-crop", "400x400+0+0",
		"-quality", "85",
		"-strip",
		thumbnailPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("errore nell'esecuzione di ImageMagick: %v", err)
	}

	return nil
}

// DeletePhoto elimina una immagine dal filesystem
func (pm *PhotoManager) DeletePhoto(filename string) error {
	filePath := filepath.Join(pm.photosDir, filename)

	// Verifica che il file esista
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file non trovato: %s", filename)
	}

	// Elimina il file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("errore nell'eliminazione del file: %v", err)
	}

	return nil
}

// isImageFile verifica se il file è un'immagine basandosi sull'estensione
func (pm *PhotoManager) isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supportedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, supportedExt := range supportedExts {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// getMimeTypeFromExtension restituisce il tipo MIME basandosi sull'estensione
func (pm *PhotoManager) getMimeTypeFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// getExtensionFromMimeType restituisce l'estensione basandosi sul MIME type
func (pm *PhotoManager) getExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg" // default
	}
}
