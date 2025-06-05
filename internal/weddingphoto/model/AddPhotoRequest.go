package model

// AddPhotoRequest rappresenta la richiesta per aggiungere una foto
type AddPhotoRequest struct {
	ImageContent string `json:"image_content" binding:"required"` // Immagine in formato base64
	ImageName    string `json:"image_name" binding:"required"`    // Nome dell'immagine
}
