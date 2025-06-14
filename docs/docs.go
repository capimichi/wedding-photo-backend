// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/photos": {
            "get": {
                "description": "Ottiene tutte le foto caricate sul server con paginazione",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "photos"
                ],
                "summary": "Recupera la lista delle foto",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Numero pagina (default: 1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Elementi per pagina (default: 10, max: 100)",
                        "name": "per_page",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.GetPhotosResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Carica una nuova foto sul server",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "photos"
                ],
                "summary": "Upload di una foto",
                "parameters": [
                    {
                        "type": "file",
                        "description": "File immagine da caricare",
                        "name": "fiimagele",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Nome personalizzato per l'immagine",
                        "name": "imageName",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.AddPhotoResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AddPhotoResponse": {
            "type": "object",
            "required": [
                "photo"
            ],
            "properties": {
                "photo": {
                    "description": "Nome della foto aggiunta",
                    "$ref": "#/definitions/model.Photo"
                }
            }
        },
        "model.ErrorResponse": {
            "type": "object",
            "required": [
                "message"
            ],
            "properties": {
                "message": {
                    "description": "Messaggio di errore",
                    "type": "string"
                }
            }
        },
        "model.GetPhotosResponse": {
            "type": "object",
            "required": [
                "page",
                "photos",
                "total_pages"
            ],
            "properties": {
                "page": {
                    "description": "Pagina corrente",
                    "type": "integer"
                },
                "photos": {
                    "description": "Lista delle foto",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Photo"
                    }
                },
                "total_pages": {
                    "description": "Numero totale di pagine",
                    "type": "integer"
                }
            }
        },
        "model.Photo": {
            "type": "object",
            "required": [
                "image_name",
                "image_url",
                "preview_url",
                "thumbnail_url"
            ],
            "properties": {
                "image_name": {
                    "description": "Nome dell'immagine",
                    "type": "string"
                },
                "image_url": {
                    "description": "URL dell'immagine",
                    "type": "string"
                },
                "preview_url": {
                    "description": "URL dell'anteprima",
                    "type": "string"
                },
                "thumbnail_url": {
                    "description": "URL del thumbnail",
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Wedding Photo Backend API",
	Description: "API per la gestione delle foto del matrimonio",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
