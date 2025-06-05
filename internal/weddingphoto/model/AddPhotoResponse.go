package model

// AddPhotoResponse rappresenta la risposta per l'aggiunta di una foto
type AddPhotoResponse struct {
	Photo Photo `json:"photo" binding:"required"` // Nome della foto aggiunta
}
