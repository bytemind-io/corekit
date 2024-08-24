/*
Copyright 2024 The corekit Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package corekit

import (
	"github.com/sashabaranov/go-openai"
	"net/http"
)

// NewInError creates a new error with a given code and message.
func NewInError(code int, err error) *openai.APIError {
	return NewError(code, err.Error())
}

// NewError creates a new error with a given code and message.
// https://platform.openai.com/docs/guides/error-codes/api-errors
func NewError(code int, message string) *openai.APIError {
	var etype string
	switch code {
	case http.StatusBadRequest:
		etype = "invalid_request_error"
	case http.StatusUnauthorized:
		etype = "authentication_error"
	case http.StatusNotFound:
		etype = "not_found_error"
	case http.StatusForbidden:
		etype = "permission_error"
	case http.StatusConflict:
		etype = "conflict_error"
	case http.StatusTooManyRequests:
		etype = "rate_limit_error"
	case http.StatusPaymentRequired:
		etype = "billing_not_active"
	case http.StatusRequestEntityTooLarge:
		etype = "request_too_large"
	case http.StatusTeapot:
		etype = "teapot_error"
	case http.StatusServiceUnavailable:
		etype = "service_unavailable_error"
	case http.StatusNotImplemented:
		etype = "api_not_implemented"
	case 529:
		etype = "overloaded_error"
	default:
		etype = "api_error"
	}

	return &openai.APIError{
		Code:           etype,
		Type:           etype,
		Message:        message,
		HTTPStatusCode: code,
	}
}
