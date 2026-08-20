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
                    "201": {"description": "Created", "schema": {"$ref": "#/definitions/httpapi.APIResponse"}},
                    "400": {"description": "Validation error", "schema": {"$ref": "#/definitions/httpapi.APIResponse"}},
                    "500": {"description": "Internal error", "schema": {"$ref": "#/definitions/httpapi.APIResponse"}}
                }
            }
        },
        "/product/{id}": {
            "patch": {
                "description": "Partial update. Omitted fields are unchanged. description and sale_price accept null.",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["product"],
                "summary": "Patch product",
                "parameters": [
                    {"type": "string", "description": "Product ID", "name": "id", "in": "path", "required": true},
                    {
                        "description": "Patch product payload",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {"$ref": "#/definitions/httpapi.PatchProductRequest"}
                    }
                ],
                "responses": {
                    "200": {"description": "OK", "schema": {"$ref": "#/definitions/httpapi.PatchAPIResponse"}},
                    "400": {"description": "Validation error", "schema": {"$ref": "#/definitions/httpapi.PatchAPIResponse"}},
                    "404": {"description": "Not found", "schema": {"$ref": "#/definitions/httpapi.PatchAPIResponse"}},
                    "500": {"description": "Internal error", "schema": {"$ref": "#/definitions/httpapi.PatchAPIResponse"}}
                }
            }
        }
    },
    "definitions": {
        "httpapi.APIResponse": {
            "type": "object",
            "properties": {
                "successful": {"type": "boolean"},
                "error_code": {"type": "string", "example": "VALIDATION_ERROR"},
                "data": {"$ref": "#/definitions/httpapi.CreateData"}
            }
        },
        "httpapi.CreateData": {
            "type": "object",
            "properties": {
                "data1": {"type": "string", "description": "Product ID"},
                "data2": {"type": "string", "description": "Product name"}
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
        "httpapi.PatchAPIResponse": {
            "type": "object",
            "properties": {
                "successful": {"type": "boolean"},
                "error_code": {"type": "string"}
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
