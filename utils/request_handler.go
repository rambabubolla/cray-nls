//
//  MIT License
//
//  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
//  Permission is hereby granted, free of charge, to any person obtaining a
//  copy of this software and associated documentation files (the "Software"),
//  to deal in the Software without restriction, including without limitation
//  the rights to use, copy, modify, merge, publish, distribute, sublicense,
//  and/or sell copies of the Software, and to permit persons to whom the
//  Software is furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included
//  in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
//  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
//  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
//  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//  OTHER DEALINGS IN THE SOFTWARE.
//
package utils

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

type ResponseError struct {
	Message string `json:"message"`
} //@name ResponseError

type ResponseOk struct {
	Message string `json:"message"`
} //@name ResponseOk

// NewRequestHandler creates a new request handler
func NewRequestHandler(logger Logger, env Env) RequestHandler {
	gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()

	instanceName := "swagger"
	if env.Environment == "development" {
		instanceName = "Internal"
	}

	engine.GET(
		"/apis/nls/openapi/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.DocExpansion("none"),
			ginSwagger.InstanceName(instanceName),
		),
	)

	return RequestHandler{Gin: engine}
}
