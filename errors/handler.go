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
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

// StatusCodeToErrorCode maps a http status code integer to an
func StatusCodeToErrorCode(statusCode int) string {
	errorCode, ok := httpStatusCodeToError[statusCode]
	if ok {
		return errorCode
	}

	return EInternal
}

// RegisterCustomStatusCode registers an error code with a status code.[注册错误码与状态码]
func RegisterCustomStatusCode(code string, statusCode int) {
	apiErrorToStatusCode[code] = statusCode
	httpStatusCodeToError[statusCode] = code
}

// ErrorCodeToStatusCode maps an influxdb error code string to a
// http status code integer.
func ErrorCodeToStatusCode(ctx context.Context, code string) int {
	// If the client disconnects early or times out then return a different
	// error than the passed in error code. Client timeouts return a 408
	// while disconnections return a non-standard Nginx HTTP 499 code.
	if err := ctx.Err(); err == context.DeadlineExceeded {
		return http.StatusRequestTimeout
	} else if err == context.Canceled {
		return 499 // https://httpstatuses.com/499
	}

	// Otherwise map internal error codes to HTTP status codes.
	statusCode, ok := apiErrorToStatusCode[code]
	if ok {
		return statusCode
	}
	return http.StatusInternalServerError
}

// apiErrorToStatusCode maps an error code to an HTTP status code.
var apiErrorToStatusCode = map[string]int{
	EEmptyValue:          http.StatusBadRequest,
	EConflict:            http.StatusConflict,
	ENotFound:            http.StatusNotFound,
	EInvalid:             http.StatusBadRequest,
	EForbidden:           http.StatusForbidden,
	ETooManyRequests:     http.StatusTooManyRequests,
	EUnauthorized:        http.StatusUnauthorized,
	EMethodNotAllowed:    http.StatusMethodNotAllowed,
	EPaymentRequired:     http.StatusPaymentRequired,
	ETooLarge:            http.StatusRequestEntityTooLarge,
	EInternal:            http.StatusInternalServerError,
	ENotImplemented:      http.StatusNotImplemented,
	EBadGateway:          http.StatusBadGateway,
	EUnprocessableEntity: http.StatusUnprocessableEntity,
	EUnavailable:         http.StatusServiceUnavailable,
}

// httpStatusCodeToError maps an HTTP status code to an error code.
var httpStatusCodeToError = map[int]string{}

// APIError is an err response body.
func init() {
	for k, v := range apiErrorToStatusCode {
		httpStatusCodeToError[v] = k
	}
}

// Err returns an error with the given code and message.
// handle RPC errors and HTTP errors simultaneously.
//
//	Example: Err(w,r, &errors.Error{
//				Code: errors.EInvalid,
//				Msg:  "model has not been found",
//			})
func Err(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	sd, ok := status.FromError(err)
	if ok {
		scode := StatusCodeToErrorCode(int(sd.Code()))

		msg := sd.Message()
		if msg == "" {
			msg = "Oops! Something went wrong."
		}

		httpx.WriteJson(w, ErrorCodeToStatusCode(r.Context(), scode), APIError{
			Code:           scode,
			Message:        msg,
			HTTPStatusCode: ErrorCodeToStatusCode(r.Context(), scode),
		})
		return
	}

	if vv, ok := err.(*APIError); ok {
		statusCode := ErrorCodeToStatusCode(r.Context(), vv.Code)
		httpx.WriteJson(w, statusCode, APIError{
			Code:           vv.Code,
			Message:        vv.Message,
			HTTPStatusCode: statusCode,
		})
		return
	}

	msg := err.Error()
	if err.Error() == "" {
		msg = "Oops! Something went wrong."
	}

	// If the error is not an APIError, return an internal error.
	httpx.WriteJson(w, http.StatusInternalServerError, APIError{
		Code:           EInternal,
		Message:        msg,
		HTTPStatusCode: ErrorCodeToStatusCode(r.Context(), EInternal),
	})
}
