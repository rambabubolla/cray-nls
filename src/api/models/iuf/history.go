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

// History
type History struct {
	ActivityState ActivityState `json:"activity_state" binding:"required" enums:"paused,in_progress,debug,blocked,wait_for_admin"` // State of activity
	SessionName   string        `json:"session_name"`                                                                              // Name of the session
	StartTime     int32         `json:"start_time"`                                                                                // Epoch timestamp
	Comment       string        `json:"comment"`                                                                                   // Comment
	Name          string        `json:"name"`                                                                                      // Comment
} //	@name	History

type ReplaceHistoryCommentRequest struct {
	Comment string `json:"comment"` // Comment
} //	@name	History.ReplaceHistoryCommentRequest

type HistoryActionRequest struct {
	StartTime int32  `json:"start_time" validate:"optional"` // Epoch timestamp
	Comment   string `json:"comment" validate:"optional"`    // Comment
} //	@name	History.HistoryActionRequest

type HistoryRunActionRequest struct {
	InputParameters InputParameters `json:"input_parameters" binding:"required"`
	SiteParameters  SiteParameters  `json:"site_parameters"`
	Comment         string          `json:"comment" validate:"optional"` // Comment
} //	@name	History.HistoryRunActionRequest
