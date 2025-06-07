package model

// AddPhotoRequest rappresenta la richiesta per aggiungere una foto
type Photo struct {
	ImageName    string `json:"image_name" binding:"required"`    // Nome dell'immagine
	ImageUrl     string `json:"image_url" binding:"required"`     // URL dell'immagine
	ThumbnailUrl string `json:"thumbnail_url" binding:"required"` // URL del thumbnail
}
