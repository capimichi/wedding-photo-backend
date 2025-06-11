#!/bin/bash

# Directory paths
MEDIA_DIR="media"
PHOTOS_DIR="$MEDIA_DIR"
PREVIEWS_DIR="$MEDIA_DIR/previews"
THUMBNAILS_DIR="$MEDIA_DIR/thumbnails"

# Redis configuration
REDIS_HOST="${REDIS_HOST:-localhost}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_PASSWORD="${REDIS_PASSWORD:-}"
QUEUE_NAME="image_processing_queue"

# Crea le directory se non esistono
mkdir -p "$PREVIEWS_DIR"
mkdir -p "$THUMBNAILS_DIR"

echo "Inizio elaborazione foto dalla coda Redis..."
echo "Connessione a Redis: $REDIS_HOST:$REDIS_PORT"

# Funzione per recuperare la prossima immagine dalla coda Redis
get_next_image_from_queue() {
    if [[ -n "$REDIS_PASSWORD" ]]; then
        redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" BRPOP "$QUEUE_NAME" 30
    else
        redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" BRPOP "$QUEUE_NAME" 30
    fi
}

# Funzione per verificare se Redis è raggiungibile
check_redis_connection() {
    if [[ -n "$REDIS_PASSWORD" ]]; then
        redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" ping >/dev/null 2>&1
    else
        redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping >/dev/null 2>&1
    fi
    return $?
}

# Verifica connessione Redis
if ! check_redis_connection; then
    echo "Errore: impossibile connettersi a Redis su $REDIS_HOST:$REDIS_PORT"
    exit 1
fi

echo "Connessione a Redis stabilita. In attesa di immagini da elaborare..."

# Loop infinito per elaborare le immagini dalla coda
processed=0
failed_attempts=0
max_failures=10

while true; do
    # Recupera la prossima immagine dalla coda (timeout 30 secondi)
    result=$(get_next_image_from_queue)
    
    # Verifica se abbiamo ricevuto un'immagine
    if [[ -n "$result" && "$result" != "(nil)" ]]; then
        # Reset counter on successful operation
        failed_attempts=0
        
        # Estrae il nome del file dal risultato di Redis
        # Il formato è: "1) queue_name\n2) filename"
        filename=$(echo "$result" | tail -n 1)
        
        if [[ -n "$filename" ]]; then
            ((processed++))
            echo "[$processed] Elaborando: $filename"
            
            # Percorso completo del file originale
            file_path="$PHOTOS_DIR/$filename"
            
            # Verifica che il file esista
            if [[ ! -f "$file_path" ]]; then
                echo "  ✗ File non trovato: $file_path"
                continue
            fi
            
            # Percorsi completi per preview e thumbnail
            preview_path="$PREVIEWS_DIR/$filename"
            thumbnail_path="$THUMBNAILS_DIR/$filename"
            
            # Crea thumbnail se non esiste
            if [[ ! -f "$thumbnail_path" ]]; then
                echo "  Creando thumbnail..."
                vipsthumbnail "$file_path" -s 400x400 -o "thumbnails/$filename"
                if [[ $? -eq 0 ]]; then
                    echo "  ✓ Thumbnail creata"
                else
                    echo "  ✗ Errore nella creazione del thumbnail"
                fi
            else
                echo "  ◦ Thumbnail già esistente"
            fi
            
            # Crea preview se non esiste
            if [[ ! -f "$preview_path" ]]; then
                echo "  Creando preview..."
                vipsthumbnail "$file_path" -s 1024x1024 -o "previews/$filename"
                if [[ $? -eq 0 ]]; then
                    echo "  ✓ Preview creata"
                else
                    echo "  ✗ Errore nella creazione della preview"
                fi
            else
                echo "  ◦ Preview già esistente"
            fi
            
            echo "  ✓ Elaborazione completata per: $filename"
            echo ""
        fi
    else
        # Nessuna immagine nella coda o errore di connessione
        ((failed_attempts++))
        echo "Tentativo fallito #$failed_attempts - Nessuna immagine in coda o errore di connessione"
        
        if [[ $failed_attempts -ge $max_failures ]]; then
            echo "Errore: raggiunti $max_failures tentativi falliti consecutivi. Uscita."
            exit 1
        fi
        
        echo "In attesa prima del prossimo tentativo..."
        sleep 5
    fi
done