package model

// ErrorResponse rappresenta la risposta per gli errori
type ErrorResponse struct {
	Message string `json:"message" binding:"required"` // Messaggio di errore
}
