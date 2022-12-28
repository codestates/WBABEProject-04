// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/menu": {
            "get": {
                "description": "등록된 메뉴 리스트를 가져온다.",
                "summary": "call GetMenuList, return ok by json.",
                "responses": {}
            },
            "post": {
                "description": "메뉴의 정보를 JSON으로 입력받아 등록한다.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "call RegisterMenu, return ok by json.",
                "responses": {}
            }
        },
        "/menu/:id": {
            "put": {
                "description": "메뉴의 아이디를 파라미터로 받고 JSON으로 수정하려는 내용을 받아 기존 메뉴의 정보를 변경할 수 있다.",
                "summary": "call UpdateMenu, return ok by json.",
                "responses": {}
            },
            "delete": {
                "description": "메뉴의 아이디를 파라미터로 받아 해당 메뉴를 삭제하는 기능",
                "summary": "call DelMenu, return ok by json.",
                "responses": {}
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
