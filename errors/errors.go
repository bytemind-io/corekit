/*
Copyright 2024 The corego Authors.

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

package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	EInternal            = "internal error"
	ENotImplemented      = "not implemented"
	EBadGateway          = "bad gateway"
	ENotFound            = "not found"
	EConflict            = "conflict"
	EInvalid             = "invalid"
	EUnprocessableEntity = "unprocessable entity"
	EEmptyValue          = "empty value"
	EUnavailable         = "unavailable"
	EForbidden           = "forbidden"
	ETooManyRequests     = "too many requests"
	EUnauthorized        = "unauthorized"
	EMethodNotAllowed    = "method not allowed"
	ETooLarge            = "request too large"
	EPaymentRequired     = "payment required"
)

// APIError is an err response body.
type APIError struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
}

// Error returns the error message.
func (e *APIError) Error() string {
	if e.Message != "" && e.Err != nil {
		var b strings.Builder
		b.WriteString(e.Message)
		b.WriteString(": ")
		b.WriteString(e.Err.Error())
		return b.String()
	} else if e.Message != "" {
		return e.Message
	} else if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("<%s>", e.Code)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *APIError) UnmarshalJSON(data []byte) (err error) {
	var rawMap map[string]json.RawMessage
	err = json.Unmarshal(data, &rawMap)
	if err != nil {
		return
	}

	err = json.Unmarshal(rawMap["message"], &e.Message)
	if err != nil {
		// If the parameter field of a function call is invalid as a JSON schema
		// refs: https://github.com/sashabaranov/go-openai/issues/381
		var messages []string
		err = json.Unmarshal(rawMap["message"], &messages)
		if err != nil {
			return
		}
		e.Message = strings.Join(messages, ", ")
	}

	if _, ok := rawMap["code"]; !ok {
		return nil
	}

	// if the api returned a number, we need to force an integer
	// since the json package defaults to float64
	var intCode string
	err = json.Unmarshal(rawMap["code"], &intCode)
	if err == nil {
		e.Code = intCode
		return nil
	}
	return json.Unmarshal(rawMap["code"], &e.Code)
}
