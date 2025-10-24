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
        "/subscriptions": {
            "post": {
                "description": "Add a new subscription to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Create a new subscription",
                "parameters": [
                    {
                        "description": "Subscription Info",
                        "name": "subscription",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.CreateSubscriptionDTO"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Subscription"
                        }
                    },
                    "400": {
                        "description": "Неверный формат JSON",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/subscriptions/summary": {
            "get": {
                "description": "Calculates the total price of subscriptions based on optional filters",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Get summary price of subscriptions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by User ID",
                        "name": "user_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Service Name",
                        "name": "service_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by start date (YYYY-MM-DD)",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by end date (YYYY-MM-DD)",
                        "name": "end_date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "integer"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid filter format",
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
        "/subscriptions/{id}": {
            "get": {
                "description": "Get details of a specific subscription",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Get a subscription by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Subscription"
                        }
                    },
                    "400": {
                        "description": "Неверный формат ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Подписка не найдена",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "put": {
                "description": "Update details of an existing subscription by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Update an existing subscription",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Subscription data to update",
                        "name": "subscription",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.UpdateSubscriptionDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Subscription not found",
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
                "description": "Delete a subscription by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "subscriptions"
                ],
                "summary": "Delete a subscription",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Subscription ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid ID format",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Subscription not found",
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
        "models.Subscription": {
            "type": "object",
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "service_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "service.CreateSubscriptionDTO": {
            "type": "object",
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "service_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "service.UpdateSubscriptionDTO": {
            "type": "object",
            "properties": {
                "end_date": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "service_name": {
                    "type": "string"
                },
                "start_date": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Subscription Service API",
	Description:      "This is a sample REST API for managing subscriptions.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
