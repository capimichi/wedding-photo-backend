package manager

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

// PhotoManager gestisce le operazioni sul filesystem per le foto
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

// createThumbnail crea un thumbnail di 200x200px
func (pm *PhotoManager) createThumbnail(originalPath, filename, contentType string) error {
    // Apre l'immagine originale
    file, err := os.Open(originalPath)
    if err != nil {
        return fmt.Errorf("errore nell'apertura del file originale: %v", err)
    }
    defer file.Close()

    // Decodifica l'immagine in base al tipo
    var img image.Image
    switch contentType {
    case "image/jpeg":
        img, err = jpeg.Decode(file)
    case "image/png":
        img, err = png.Decode(file)
    case "image/gif":
        img, err = gif.Decode(file)
    case "image/webp":
        img, err = webp.Decode(file)
    default:
        img, _, err = image.Decode(file)
    }

    if err != nil {
        return fmt.Errorf("errore nella decodifica dell'immagine: %v", err)
    }

    // Crea un'immagine thumbnail di 200x200
    thumbnail := image.NewRGBA(image.Rect(0, 0, 200, 200))

    // Calcola il rettangolo di crop per mantenere le proporzioni
    srcBounds := img.Bounds()
    srcWidth := srcBounds.Dx()
    srcHeight := srcBounds.Dy()

    // Calcola il rapporto di aspetto
    srcAspect := float64(srcWidth) / float64(srcHeight)
    dstAspect := 1.0 // 200x200 è quadrato

    var cropRect image.Rectangle
    if srcAspect > dstAspect {
        // Immagine più larga: crop orizzontalmente
        newWidth := int(float64(srcHeight) * dstAspect)
        offset := (srcWidth - newWidth) / 2
        cropRect = image.Rect(srcBounds.Min.X+offset, srcBounds.Min.Y, srcBounds.Min.X+offset+newWidth, srcBounds.Max.Y)
    } else {
        // Immagine più alta: crop verticalmente
        newHeight := int(float64(srcWidth) / dstAspect)
        offset := (srcHeight - newHeight) / 2
        cropRect = image.Rect(srcBounds.Min.X, srcBounds.Min.Y+offset, srcBounds.Max.X, srcBounds.Min.Y+offset+newHeight)
    }

    // Usa BiLinear per una qualità migliore nei thumbnail
    draw.BiLinear.Scale(thumbnail, thumbnail.Bounds(), img, cropRect, draw.Over, nil)

    // Salva il thumbnail
    thumbnailPath := filepath.Join(pm.thumbnailsDir, filename)
    thumbnailFile, err := os.Create(thumbnailPath)
    if err != nil {
        return fmt.Errorf("errore nella creazione del file thumbnail: %v", err)
    }
    defer thumbnailFile.Close()

    // Aumenta la qualità JPEG per i thumbnail
    return jpeg.Encode(thumbnailFile, thumbnail, &jpeg.Options{Quality: 95})
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
