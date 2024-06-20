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
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
)

// ErrBody is an err response body.
type ErrBody struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

// Err returns an error with the given code and message.
// handle RPC errors and HTTP errors simultaneously.
func Err(w http.ResponseWriter, err error) {
	stu, ok := status.FromError(err)
	if ok {
		httpx.WriteJson(w, int(stu.Code()), ErrBody{
			Code: fmt.Sprintf("%d", stu.Code()),
			Msg:  stu.Message(),
		})
		return
	}

	if vv, ok := err.(*errors.CodeMsg); ok {
		httpx.WriteJson(w, vv.Code, ErrBody{
			Code: fmt.Sprintf("%d", vv.Code),
			Msg:  vv.Msg,
		})
		return
	}

	httpx.WriteJson(w, http.StatusInternalServerError, ErrBody{
		Code: "500",
		Msg:  err.Error(),
	})
}
