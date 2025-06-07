package model

// GetPhotosResponse rappresenta la risposta per il recupero delle foto con paginazione
type GetPhotosResponse struct {
	Photos     []Photo `json:"photos" binding:"required"`      // Lista delle foto
	Page       int     `json:"page" binding:"required"`        // Pagina corrente
	TotalPages int     `json:"total_pages" binding:"required"` // Numero totale di pagine
}
