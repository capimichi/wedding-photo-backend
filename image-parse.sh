#!/bin/bash

# Directory paths
MEDIA_DIR="media"
PHOTOS_DIR="$MEDIA_DIR"
PREVIEWS_DIR="$MEDIA_DIR/previews"
THUMBNAILS_DIR="$MEDIA_DIR/thumbnails"

# Crea le directory se non esistono
mkdir -p "$PREVIEWS_DIR"
mkdir -p "$THUMBNAILS_DIR"

# Estensioni supportate (case insensitive)
EXTENSIONS=("jpg" "jpeg" "png" "gif" "webp")

echo "Inizio elaborazione foto..."

# Funzione per verificare se un file è un'immagine
is_image_file() {
    local file="$1"
    local ext=$(echo "${file##*.}" | tr '[:upper:]' '[:lower:]')
    
    for supported_ext in "${EXTENSIONS[@]}"; do
        if [[ "$ext" == "$supported_ext" ]]; then
            return 0
        fi
    done
    return 1
}

# Conta il numero totale di foto da elaborare
total_photos=0
for file in "$PHOTOS_DIR"/*; do
    if [[ -f "$file" ]] && is_image_file "$file"; then
        ((total_photos++))
    fi
done

echo "Trovate $total_photos foto da elaborare"

# Contatore per il progresso
processed=0

# Itera attraverso tutti i file nella directory media
for file in "$PHOTOS_DIR"/*; do
    # Verifica che sia un file (non una directory)
    if [[ -f "$file" ]]; then
        filename=$(basename "$file")
        
        # Verifica che sia un'immagine
        if is_image_file "$filename"; then
            ((processed++))
            echo "[$processed/$total_photos] Elaborando: $filename"
            
            # Percorsi completi per preview e thumbnail
            preview_path="$PREVIEWS_DIR/$filename"
            thumbnail_path="$THUMBNAILS_DIR/$filename"
            
            # Crea preview se non esiste
            if [[ ! -f "$preview_path" ]]; then
                echo "  Creando preview..."
                magick "$file" -auto-orient -resize 1024x1024\> -quality 85 "$preview_path"
                if [[ $? -eq 0 ]]; then
                    echo "  ✓ Preview creata"
                else
                    echo "  ✗ Errore nella creazione della preview"
                fi
            else
                echo "  ◦ Preview già esistente"
            fi
            
            # Crea thumbnail se non esiste
            if [[ ! -f "$thumbnail_path" ]]; then
                echo "  Creando thumbnail..."
                magick "$file" -auto-orient -resize 400x400^ -gravity center -extent 400x400 -quality 85 "$thumbnail_path"
                if [[ $? -eq 0 ]]; then
                    echo "  ✓ Thumbnail creata"
                else
                    echo "  ✗ Errore nella creazione del thumbnail"
                fi
            else
                echo "  ◦ Thumbnail già esistente"
            fi
            
            echo ""
        fi
    fi
done

echo "Elaborazione completata!"
echo "Foto elaborate: $processed"