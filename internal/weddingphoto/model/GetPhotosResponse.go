package model

// AddPhotoResponse rappresenta la risposta per l'aggiunta di una foto
type GetPhotosResponse struct {
	Photos []Photo `json:"photos" binding:"required"` // Nome della foto aggiunta
}
