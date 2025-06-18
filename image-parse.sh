#!/bin/bash

# Directory paths
MEDIA_DIR="media"
PHOTOS_DIR="$MEDIA_DIR"
PREVIEWS_DIR="$MEDIA_DIR/previews"
THUMBNAILS_DIR="$MEDIA_DIR/thumbnails"

# Crea le directory se non esistono
mkdir -p "$PREVIEWS_DIR"
mkdir -p "$THUMBNAILS_DIR"

while true; do
    found_files=0

    # Trova i file creati/modificati negli ultimi 240 minuti
    while IFS= read -r file; do
        filename=$(basename "$file")
        previews_path="$PREVIEWS_DIR/$filename"
        thumbnails_path="$THUMBNAILS_DIR/$filename"
        file_path="$file"

        # Crea preview se non esiste
        if [[ ! -f "$previews_path" ]]; then
            echo "Creando preview per: $filename"
            vipsthumbnail "$file_path" -s 1024x1024 -o "$previews_path"
            if [[ $? -eq 0 ]]; then
                echo "Preview creata con successo: $filename"
            else
                echo "Errore nella creazione della preview per: $filename"
            fi
            found_files=1
        fi

        # Crea thumbnail se non esiste
        if [[ ! -f "$thumbnails_path" ]]; then
            echo "Creando thumbnail per: $filename"
            vipsthumbnail "$previews_path" -s 400x400 -o "$thumbnails_path"
            if [[ $? -eq 0 ]]; then
                echo "Thumbnail creata con successo: $filename"
            else
                echo "Errore nella creazione della thumbnail per: $filename"
            fi
            found_files=1
        fi
    done < <(find "$PHOTOS_DIR" -type f -mmin -240)

    # Se non sono stati trovati file da elaborare, aspetta 5 secondi
    if [[ $found_files -eq 0 ]]; then
        echo "Nessun file da elaborare. Attendo 5 secondi..."
        sleep 5
    fi
done
