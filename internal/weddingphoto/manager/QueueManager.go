package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	IMAGE_PROCESSING_QUEUE = "image_processing_queue"
)

// QueueManager gestisce la comunicazione con Redis per la coda di elaborazione immagini
type QueueManager struct {
	client *redis.Client
	ctx    context.Context
}

// NewQueueManager crea una nuova istanza del manager
func NewQueueManager(redisAddr, redisPassword string, redisDB int) *QueueManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})

	return &QueueManager{
		client: rdb,
		ctx:    context.Background(),
	}
}

// AddImageToQueue aggiunge un'immagine alla coda di elaborazione
func (qm *QueueManager) AddImageToQueue(imageName string) error {
	err := qm.client.LPush(qm.ctx, IMAGE_PROCESSING_QUEUE, imageName).Err()
	if err != nil {
		return fmt.Errorf("errore nell'aggiunta dell'immagine alla coda: %v", err)
	}
	return nil
}

// GetNextImageFromQueue recupera la prossima immagine dalla coda (operazione bloccante)
func (qm *QueueManager) GetNextImageFromQueue(timeout time.Duration) (string, error) {
	result, err := qm.client.BRPop(qm.ctx, timeout, IMAGE_PROCESSING_QUEUE).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil // Nessun elemento nella coda
		}
		return "", fmt.Errorf("errore nel recupero dell'immagine dalla coda: %v", err)
	}

	if len(result) < 2 {
		return "", fmt.Errorf("risposta Redis malformata")
	}

	return result[1], nil // result[0] è il nome della coda, result[1] è il valore
}

// GetQueueLength restituisce il numero di elementi nella coda
func (qm *QueueManager) GetQueueLength() (int64, error) {
	length, err := qm.client.LLen(qm.ctx, IMAGE_PROCESSING_QUEUE).Result()
	if err != nil {
		return 0, fmt.Errorf("errore nel recupero della lunghezza della coda: %v", err)
	}
	return length, nil
}

// TestConnection testa la connessione a Redis
func (qm *QueueManager) TestConnection() error {
	_, err := qm.client.Ping(qm.ctx).Result()
	if err != nil {
		return fmt.Errorf("errore nella connessione a Redis: %v", err)
	}
	return nil
}

// Close chiude la connessione Redis
func (qm *QueueManager) Close() error {
	return qm.client.Close()
}
