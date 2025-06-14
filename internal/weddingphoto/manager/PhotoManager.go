package manager

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

// PhotoManager gestisce le operazioni sulfilesystem per le foto
type PhotoManager struct {
	photosDir     string
	thumbnailsDir string
	previewsDir   string
}

// NewPhotoManager crea una nuova istanza del manager
func NewPhotoManager(photosDir string) *PhotoManager {
	thumbnailsDir := filepath.Join(photosDir, "thumbnails")
	previewsDir := filepath.Join(photosDir, "previews")

	// Crea le directory se non esistono
	if err := os.MkdirAll(photosDir, 0755); err != nil {
		fmt.Printf("Errore nella creazione della directory %s: %v\n", photosDir, err)
	}
	if err := os.MkdirAll(thumbnailsDir, 0755); err != nil {
		fmt.Printf("Errore nella creazione della directory thumbnails %s: %v\n", thumbnailsDir, err)
	}
	if err := os.MkdirAll(previewsDir, 0755); err != nil {
		fmt.Printf("Errore nella creazione della directory previews %s: %v\n", previewsDir, err)
	}

	return &PhotoManager{
		photosDir:     photosDir,
		thumbnailsDir: thumbnailsDir,
		previewsDir:   previewsDir,
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

	// Crea la preview
	if err := pm.createPreview(filePath, filename, contentType); err != nil {
		fmt.Printf("Errore nella creazione della preview per %s: %v\n", filename, err)
		// Non restituiamo errore per la preview, continuiamo
	}

	return filename, nil
}

// createThumbnail crea un thumbnail di 400x400px usando la libreria imaging
func (pm *PhotoManager) createThumbnail(originalPath, filename, contentType string) error {
	return nil
	thumbnailPath := filepath.Join(pm.thumbnailsDir, filename)

	// Apre l'immagine originale
	src, err := imaging.Open(originalPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("errore nell'apertura dell'immagine: %v", err)
	}

	// Crea il thumbnail 400x400 con crop al centro
	thumbnail := imaging.Fill(src, 400, 400, imaging.Center, imaging.Lanczos)

	// Salva il thumbnail con qualità JPEG 85
	err = imaging.Save(thumbnail, thumbnailPath, imaging.JPEGQuality(85))
	if err != nil {
		return fmt.Errorf("errore nel salvataggio del thumbnail: %v", err)
	}

	return nil
}

// createPreview crea una preview con dimensioni massime 1024x1024 mantenendo le proporzioni
func (pm *PhotoManager) createPreview(originalPath, filename, contentType string) error {
	return nil
	previewPath := filepath.Join(pm.previewsDir, filename)

	// Apre l'immagine originale
	src, err := imaging.Open(originalPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("errore nell'apertura dell'immagine: %v", err)
	}

	// Ridimensiona mantenendo le proporzioni con dimensioni massime 1024x1024
	preview := imaging.Fit(src, 1024, 1024, imaging.Lanczos)

	// Salva la preview con qualità JPEG 85
	err = imaging.Save(preview, previewPath, imaging.JPEGQuality(85))
	if err != nil {
		return fmt.Errorf("errore nel salvataggio della preview: %v", err)
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

// ThumbnailExists verifica se il thumbnail di un'immagine esiste
func (pm *PhotoManager) ThumbnailExists(filename string) bool {
	thumbnailPath := filepath.Join(pm.thumbnailsDir, filename)
	_, err := os.Stat(thumbnailPath)
	return !os.IsNotExist(err)
}

// PreviewExists verifica se la preview di un'immagine esiste
func (pm *PhotoManager) PreviewExists(filename string) bool {
	previewPath := filepath.Join(pm.previewsDir, filename)
	_, err := os.Stat(previewPath)
	return !os.IsNotExist(err)
}

// DetectMimeTypeFromBytes rileva il MIME type reale leggendo i magic bytes
func (pm *PhotoManager) DetectMimeTypeFromBytes(reader io.Reader) (string, io.Reader, error) {
	// Legge i primi 512 bytes per il rilevamento del MIME type
	buffer := make([]byte, 512)
	n, err := reader.Read(buffer)
	if err != nil && err != io.EOF {
		return "", nil, fmt.Errorf("errore nella lettura dei bytes: %v", err)
	}

	// Crea un nuovo reader che include i bytes letti + il resto del file originale
	newReader := io.MultiReader(bytes.NewReader(buffer[:n]), reader)

	// Rileva il MIME type dai magic bytes
	mimeType := pm.detectMimeFromMagicBytes(buffer[:n])

	return mimeType, newReader, nil
}

// detectMimeFromMagicBytes rileva il MIME type dai magic bytes
func (pm *PhotoManager) detectMimeFromMagicBytes(data []byte) string {
	if len(data) < 8 {
		return ""
	}

	// JPEG
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}

	// PNG
	if len(data) >= 8 && bytes.Equal(data[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
		return "image/png"
	}

	// GIF87a or GIF89a
	if len(data) >= 6 && (bytes.Equal(data[:6], []byte("GIF87a")) || bytes.Equal(data[:6], []byte("GIF89a"))) {
		return "image/gif"
	}

	// WebP
	if len(data) >= 12 && bytes.Equal(data[:4], []byte("RIFF")) && bytes.Equal(data[8:12], []byte("WEBP")) {
		return "image/webp"
	}

	return ""
}

// IsValidImageMimeType verifica se il MIME type è di un'immagine supportata
func (pm *PhotoManager) IsValidImageMimeType(mimeType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if mimeType == validType {
			return true
		}
	}
	return false
}
