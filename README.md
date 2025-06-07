# Wedding Photo Backend

Backend per il progetto di gestione foto matrimonio, sviluppato in Go con framework Gin.

## Funzionalita

- Upload di foto tramite REST API
- Recupero della lista delle foto caricate
- Supporto per immagini JPEG, PNG, GIF e WebP
- Validazione dei file caricati
- CORS abilitato per integrazioni frontend

## API Endpoints

### POST /api/photos
Carica una nuova foto.

**Parametri:**
- `photo` (file): File immagine da caricare (multipart/form-data)

**Risposta di successo (200):**
```json
{
  "success": true,
  "message": "Foto caricata con successo",
  "photo": {
    "id": "1672531200",
    "filename": "esempio.jpg",
    "path": "uploads/1672531200_esempio.jpg",
    "size": 1024567,
    "mime_type": "image/jpeg",
    "upload_time": "2025-06-03T10:30:00Z"
  }
}
```

### GET /api/photos
Recupera la lista di tutte le foto caricate.

**Risposta di successo (200):**
```json
{
  "success": true,
  "photos": [
    {
      "id": "1672531200",
      "filename": "esempio.jpg",
      "path": "uploads/1672531200_esempio.jpg",
      "size": 1024567,
      "mime_type": "image/jpeg",
      "upload_time": "2025-06-03T10:30:00Z"
    }
  ],
  "count": 1
}
```

## Avvio del server

```bash
# Installa le dipendenze
go mod tidy

# Avvia il server
go run main.go
```

Il server sarà disponibile su `http://localhost:8080`

## Struttura del progetto

```
wedding-photo-backend/
├── main.go                                 # Entry point dell'applicazione
├── go.mod                                  # Gestione dipendenze Go
├── internal/
│   └── weddingphoto/
│       └── controller/
│           └── PhotoController.go          # Controller per gestione foto
└── uploads/                               # Directory per file caricati (creata automaticamente)
```

## Tecnologie utilizzate

- **Go 1.21+**
- **Gin** - Web framework
- **Standard library** - Per gestione file e HTTP
