// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "VictorNM",
            "url": "https://github.com/VictorNM/"
        },
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users/profile": {
            "get": {
                "description": "Get profile by user_id in token,",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get current sign-inned user's profile",
                "responses": {
                    "200": {
                        "description": "Get profile successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.BaseResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/user.ProfileDTO"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users/register": {
            "post": {
                "description": "Sign in using email and password",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Basic sign in using email, password",
                "parameters": [
                    {
                        "description": "Register new user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.RegisterMutation"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Register successfully",
                        "schema": {
                            "$ref": "#/definitions/api.BaseResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.BaseResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/api.Error"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/users/sign-in": {
            "post": {
                "description": "Sign in using email and password",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Basic sign in using email, password",
                "responses": {
                    "200": {
                        "description": "Sign in successfully",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.BaseResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/api.authToken"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Not authenticated",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/api.BaseResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "errors": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/api.Error"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.BaseResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "errors": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.Error"
                    }
                }
            }
        },
        "api.Error": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "api.authToken": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "user.ProfileDTO": {
            "type": "object",
            "properties": {
                "country": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "language": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "year_of_birth": {
                    "type": "integer"
                }
            }
        },
        "user.RegisterMutation": {
            "type": "object",
            "required": [
                "email",
                "full_name",
                "password",
                "password_confirmation",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "full_name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "password_confirmation": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
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
	Host:        "localhost:8080",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "ES API",
	Description: "",
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
	swag.Register(swag.Name, &s{})
}
