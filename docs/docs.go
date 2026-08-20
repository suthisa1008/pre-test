package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/product": {
            "post": {
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["product"],
                "summary": "Create product",
                "parameters": [
                    {
                        "description": "Create product payload",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/httpapi.CreateProductRequest"}
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {"$ref": "#/definitions/httpapi.CreateSuccessResponse"},
                        "examples": {
                            "application/json": {
                                "successful": true,
                                "error_code": "SUCCESS",
                                "data": {
                                    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
                                    "name": "Latte",
                                    "description": "hot milk coffee",
                                    "sale_price": 80,
                                    "price": 100,
                                    "created_at": "2026-08-20T10:00:00Z",
                                    "updated_at": "2026-08-20T10:00:00Z"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid JSON or domain validation failed (empty name, negative price, sale_price greater than price)",
                        "schema": {"$ref": "#/definitions/httpapi.CreateValidationErrorResponse"},
                        "examples": {
                            "application/json": {
                                "successful": false,
                                "error_code": "VALIDATION_ERROR",
                                "data": null
                            }
                        }
                    },
                    "500": {
                        "description": "Unexpected persistence or internal failure",
                        "schema": {"$ref": "#/definitions/httpapi.CreateInternalErrorResponse"},
                        "examples": {
                            "application/json": {
                                "successful": false,
                                "error_code": "INTERNAL_ERROR",
                                "data": null
                            }
                        }
                    }
                }
            }
        },
        "/product/{id}": {
            "patch": {
                "description": "Partial update. Omitted fields are unchanged. description and sale_price accept null. Response has no data field.",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["product"],
                "summary": "Patch product",
                "parameters": [
                    {"type": "string", "description": "Product ID", "name": "id", "in": "path", "required": true},
                    {
                        "description": "Patch product payload. Send only fields to change.",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/httpapi.PatchProductRequest"}
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Updated",
                        "schema": {"$ref": "#/definitions/httpapi.PatchSuccessResponse"},
                        "examples": {
                            "application/json": {
                                "successful": true,
                                "error_code": "SUCCESS"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid JSON, empty body, name/price set to null, or domain validation failed",
                        "schema": {"$ref": "#/definitions/httpapi.PatchValidationErrorResponse"},
                        "examples": {
                            "application/json": {
                                "successful": false,
                                "error_code": "VALIDATION_ERROR"
                            }
                        }
                    },
                    "404": {
                        "description": "Product id does not exist",
                        "schema": {"$ref": "#/definitions/httpapi.PatchNotFoundErrorResponse"},
                        "examples": {
                            "application/json": {
                                "successful": false,
                                "error_code": "NOT_FOUND"
                            }
                        }
                    },
                    "500": {
                        "description": "Unexpected persistence or internal failure",
                        "schema": {"$ref": "#/definitions/httpapi.PatchInternalErrorResponse"},
                        "examples": {
                            "application/json": {
                                "successful": false,
                                "error_code": "INTERNAL_ERROR"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "httpapi.ProductResponse": {
            "type": "object",
            "required": ["id", "name", "price", "created_at", "updated_at"],
            "properties": {
                "id": {"type": "string", "format": "uuid"},
                "name": {"type": "string"},
                "description": {"type": "string", "x-nullable": true},
                "sale_price": {"type": "number", "x-nullable": true},
                "price": {"type": "number"},
                "created_at": {"type": "string", "format": "date-time"},
                "updated_at": {"type": "string", "format": "date-time"}
            }
        },
        "httpapi.CreateSuccessResponse": {
            "type": "object",
            "required": ["successful", "error_code", "data"],
            "properties": {
                "successful": {"type": "boolean", "example": true},
                "error_code": {"type": "string", "enum": ["SUCCESS"], "example": "SUCCESS"},
                "data": {"$ref": "#/definitions/httpapi.ProductResponse"}
            }
        },
        "httpapi.CreateValidationErrorResponse": {
            "type": "object",
            "required": ["successful", "error_code", "data"],
            "properties": {
                "successful": {"type": "boolean", "example": false},
                "error_code": {"type": "string", "enum": ["VALIDATION_ERROR"], "example": "VALIDATION_ERROR"},
                "data": {"description": "Always null on error", "x-nullable": true, "example": null}
            }
        },
        "httpapi.CreateInternalErrorResponse": {
            "type": "object",
            "required": ["successful", "error_code", "data"],
            "properties": {
                "successful": {"type": "boolean", "example": false},
                "error_code": {"type": "string", "enum": ["INTERNAL_ERROR"], "example": "INTERNAL_ERROR"},
                "data": {"description": "Always null on error", "x-nullable": true, "example": null}
            }
        },
        "httpapi.CreateProductRequest": {
            "type": "object",
            "required": ["name", "price"],
            "properties": {
                "name": {"type": "string"},
                "description": {"type": "string", "x-nullable": true},
                "sale_price": {"type": "number", "x-nullable": true},
                "price": {"type": "number"}
            }
        },
        "httpapi.PatchSuccessResponse": {
            "type": "object",
            "required": ["successful", "error_code"],
            "properties": {
                "successful": {"type": "boolean", "example": true},
                "error_code": {"type": "string", "enum": ["SUCCESS"], "example": "SUCCESS"}
            }
        },
        "httpapi.PatchValidationErrorResponse": {
            "type": "object",
            "required": ["successful", "error_code"],
            "properties": {
                "successful": {"type": "boolean", "example": false},
                "error_code": {"type": "string", "enum": ["VALIDATION_ERROR"], "example": "VALIDATION_ERROR"}
            }
        },
        "httpapi.PatchNotFoundErrorResponse": {
            "type": "object",
            "required": ["successful", "error_code"],
            "properties": {
                "successful": {"type": "boolean", "example": false},
                "error_code": {"type": "string", "enum": ["NOT_FOUND"], "example": "NOT_FOUND"}
            }
        },
        "httpapi.PatchInternalErrorResponse": {
            "type": "object",
            "required": ["successful", "error_code"],
            "properties": {
                "successful": {"type": "boolean", "example": false},
                "error_code": {"type": "string", "enum": ["INTERNAL_ERROR"], "example": "INTERNAL_ERROR"}
            }
        },
        "httpapi.PatchProductRequest": {
            "type": "object",
            "properties": {
                "name": {"type": "string"},
                "description": {"type": "string", "x-nullable": true},
                "sale_price": {"type": "number", "x-nullable": true},
                "price": {"type": "number"}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it.
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Product API",
	Description:      "Product REST API with PostgreSQL, Clean Architecture, and partial PATCH.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
