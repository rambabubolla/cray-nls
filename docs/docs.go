// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
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
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/ncns/{hostname}/reboot": {
            "post": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "summary": "End to end reboot of a single ncn",
                "parameters": [
                    {
                        "type": "string",
                        "description": "hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v1/ncns/{hostname}/rebuild": {
            "post": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN"
                ],
                "summary": "End to end rebuild of a single ncn",
                "parameters": [
                    {
                        "type": "string",
                        "description": "hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v1/workflows": {
            "get": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin",
                            "read"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow"
                ],
                "summary": "Get status of a ncn workflow",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "collectionFormat": "csv",
                        "description": "workflow ids",
                        "name": "workflow_ids",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v1/workflows/{name}": {
            "delete": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin",
                            "read"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow"
                ],
                "summary": "Delete a ncn workflow",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of workflow",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v1/workflows/{name}/rerun": {
            "put": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin",
                            "read"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow"
                ],
                "summary": "Rerun a workflow, all steps will run",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of workflow",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v1/workflows/{name}/retry": {
            "put": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin",
                            "read"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Workflow"
                ],
                "summary": "Retry a failed ncn workflow, skip passed steps",
                "parameters": [
                    {
                        "type": "string",
                        "description": "name of workflow",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v2/ncn": {
            "post": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN v2"
                ],
                "summary": "Add a ncn",
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v2/ncns/reboot": {
            "post": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN v2"
                ],
                "summary": "End to end rolling reboot request",
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v2/ncns/rebuild": {
            "post": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN v2"
                ],
                "summary": "End to end rolling rebuild request",
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        },
        "/v2/ncns/{hostname}": {
            "delete": {
                "security": [
                    {
                        "OAuth2Application": [
                            "admin"
                        ]
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "NCN v2"
                ],
                "summary": "Remove a ncn",
                "parameters": [
                    {
                        "type": "string",
                        "description": "hostname",
                        "name": "hostname",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "501": {
                        "description": "Not Implemented"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": "                             Grants read and write access to administrative information",
                "read": "                              Grants read access"
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/apis/nls",
	Schemes:          []string{},
	Title:            "NCN Lifecycle Management API",
	Description:      "This doc descibes REST API for ncn lifecycle management. Note that in this version, we only provide APIs for individual operation. A full end to end lifecycle management API is out of scope in Phase I\n\n---\n\n## Argo workflow Demo\n\n---\n\n[API Doc](swagger.md)\n",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
