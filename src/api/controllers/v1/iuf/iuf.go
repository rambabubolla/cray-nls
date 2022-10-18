/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package iuf

import (
	"github.com/Cray-HPE/cray-nls/src/api/services"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/gin-gonic/gin"
)

// IufController data type
type IufController struct {
	workflowService services.WorkflowService
	logger          utils.Logger
}

// NewIufController creates new Ncn controller
func NewIufController(workflowService services.WorkflowService, logger utils.Logger) IufController {
	return IufController{
		workflowService: workflowService,
		logger:          logger,
	}
}

// IufCreateInstall
// @Summary  This is a summary
// @Param    hostname  path  string  true  "hostname"
// @Tags     This_IS_A_TAG
// @Accept   json
// @Produce  json
// @Failure  501  "Not Implemented"
// @Router   /iuf/v1/install [post]
func (u IufController) IufCreateInstall(c *gin.Context) {
	c.JSON(501, "not implemented")
}
