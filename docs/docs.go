// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate_swagger = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
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
        "/inventory": {
            "get": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Inventory"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Inventory"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ApiResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "Create Inventory",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Inventory"
                ],
                "parameters": [
                    {
                        "description": "Add Inventory",
                        "name": "Inventory",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Inventory"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ApiResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ApiResponse"
                        }
                    }
                }
            }
        },
        "/inventory/{id}": {
            "get": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Inventory"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Inventory ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Inventory"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ApiResponse"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "get struct array by ID",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "parameters": [
                    {
                        "description": "Add account",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponseFailed"
                        }
                    }
                }
            }
        },
        "/shop": {
            "get": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "Access to Shop",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Shop"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ShopInfo"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ApiResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "Create Shop",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Shop"
                ],
                "parameters": [
                    {
                        "description": "Add Shop",
                        "name": "Shop",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Shop"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Shop"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSuccessWithId"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "description": "For User Register Application",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Register An Account",
                "parameters": [
                    {
                        "description": "Add account",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ResponseSuccessWithId"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.AuthResponseFailed"
                        }
                    }
                }
            }
        },
        "/select-shop": {
            "post": {
                "security": [
                    {
                        "AccessToken": []
                    }
                ],
                "description": "Access to Shop",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "parameters": [
                    {
                        "description": "Shop",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ShopSelectRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ApiResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.ApiResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ApiResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "id": {},
                "message": {
                    "type": "string"
                },
                "pagination": {},
                "success": {
                    "type": "boolean"
                }
            }
        },
        "models.AuthResponse": {
            "type": "object",
            "properties": {
                "success": {
                    "type": "boolean"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.AuthResponseFailed": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "models.Barcode": {
            "type": "object",
            "properties": {
                "barcode": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "unit": {
                    "type": "string"
                }
            }
        },
        "models.Inventory": {
            "type": "object",
            "properties": {
                "activated": {
                    "description": "เปิดใช้งานอยู่",
                    "type": "boolean"
                },
                "barcodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Barcode"
                    }
                },
                "categoryGuid": {
                    "description": "Guid กลุ่มสินค้า",
                    "type": "string"
                },
                "description1": {
                    "description": "รายละเอียดภาษาไทย",
                    "type": "string"
                },
                "description2": {
                    "type": "string"
                },
                "description3": {
                    "type": "string"
                },
                "description4": {
                    "type": "string"
                },
                "description5": {
                    "type": "string"
                },
                "guidFixed": {
                    "description": "Guid สินค้า",
                    "type": "string"
                },
                "haveImage": {
                    "description": "มีรูปสินค้า",
                    "type": "boolean"
                },
                "id": {
                    "type": "string"
                },
                "itemSku": {
                    "type": "string"
                },
                "lineNumber": {
                    "description": "บรรทัดที่ (เอาไว้เรียงลำดับ)",
                    "type": "integer"
                },
                "shopId": {
                    "description": "รหัสร้าน",
                    "type": "string"
                },
                "name1": {
                    "description": "ชื่อภาษาไทย",
                    "type": "string"
                },
                "name2": {
                    "type": "string"
                },
                "name3": {
                    "type": "string"
                },
                "name4": {
                    "type": "string"
                },
                "name5": {
                    "type": "string"
                },
                "price": {
                    "description": "ราคาพื้นฐาน (กรณีไม่มีตารางราคา และโปรโมชั่น)",
                    "type": "number"
                },
                "recommended": {
                    "description": "สินค้าแนะนำ",
                    "type": "boolean"
                },
                "unitList": {
                    "description": "กรณีหลายหน่วยนับ ตารางหน่วบนับ",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.InventoryUnit"
                    }
                }
            }
        },
        "models.InventoryUnit": {
            "type": "object",
            "properties": {
                "activated": {
                    "description": "เปิดใช้งานอยู่",
                    "type": "boolean"
                },
                "divisor": {
                    "description": "ตัวหาร",
                    "type": "number"
                },
                "minuend": {
                    "description": "ตัวตั้ง",
                    "type": "number"
                },
                "unitGuid": {
                    "description": "Guid หน่วยนับ",
                    "type": "string"
                },
                "unitName": {
                    "description": "ชื่อหน่วยนับ",
                    "type": "string"
                }
            }
        },
        "models.Shop": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "default": "-"
                },
                "name1": {
                    "type": "string"
                }
            }
        },
        "models.ShopInfo": {
            "type": "object",
            "properties": {
                "guidFixed": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name1": {
                    "type": "string"
                }
            }
        },
        "models.ShopSelectRequest": {
            "type": "object",
            "properties": {
                "shopId": {
                    "type": "string"
                }
            }
        },
        "models.ResponseSuccessWithId": {
            "type": "object",
            "properties": {
                "id": {},
                "success": {
                    "type": "boolean"
                }
            }
        },
        "models.UserRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "AccessToken": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo_swagger holds exported Swagger Info so clients can modify it
var SwaggerInfo_swagger = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "SML Cloud Platform API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo_swagger.InstanceName(), SwaggerInfo_swagger)
}
