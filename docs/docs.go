// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Returns a list of songs with filtering and pagination capabilities",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get library data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by Group",
                        "name": "group_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Songname",
                        "name": "song_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by ReleaseDate",
                        "name": "release_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Text",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Link",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page Number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of elements on one page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Songs list",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Song"
                            }
                        }
                    },
                    "400": {
                        "description": "Wrong request params",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Adding new song to DB",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Add new song",
                "parameters": [
                    {
                        "description": "New song data",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SongRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Song was added successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Wrong request params",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/songs/{songID}": {
            "put": {
                "description": "Updates song data in the database by song ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Update song",
                "parameters": [
                    {
                        "type": "string",
                        "description": "song ID",
                        "name": "songID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated song details",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Updated song",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    },
                    "400": {
                        "description": "Wrong request data",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Song is not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Deletes a song from the database by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Delete song",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Song ID",
                        "name": "songID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song was deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Song is not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/songs/{songID}/lyrics": {
            "get": {
                "description": "Returns lyrics of a song with pagynation capability",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get song lyrics",
                "parameters": [
                    {
                        "type": "string",
                        "description": "song ID",
                        "name": "songID",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of verses per page",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Song text",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Wrong request params",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Song is not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Song": {
            "description": "Represents a song in the database.",
            "type": "object",
            "properties": {
                "group": {
                    "description": "@Description The name of the music group.",
                    "type": "string"
                },
                "id": {
                    "description": "@Description The unique identifier of the song.",
                    "type": "integer"
                },
                "link": {
                    "description": "@Description The external link to the song (e.g.,\"https://www.youtube.com/watch?v=Xsp3_a-PMTw\").",
                    "type": "string"
                },
                "release_date": {
                    "description": "@Description The release date of the song (e.g., \"16.07.2006\").",
                    "type": "string"
                },
                "song": {
                    "description": "@Description The name of the song.",
                    "type": "string"
                },
                "text": {
                    "description": "@Description The lyrics or text of the song.",
                    "type": "string"
                }
            }
        },
        "models.SongRequest": {
            "description": "Represents the request payload for adding a new song.",
            "type": "object",
            "properties": {
                "group": {
                    "description": "@Description The name of the music group.",
                    "type": "string"
                },
                "song": {
                    "description": "@Description The name of the song.",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "",
	Schemes:          []string{"http"},
	Title:            "Music Library API",
	Description:      "API for music library management",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
