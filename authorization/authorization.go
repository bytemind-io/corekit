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

package authorization

import (
	"net/http"
	"strings"
)

// IsBasicAuth checks if the request is using basic auth.
func IsBasicAuth(r *http.Request) bool {
	_, _, ok := r.BasicAuth()
	return ok
}

// IsAuthorizationToken checks if the request is using authorization token.
func IsAuthorizationToken(r *http.Request) bool {
	return r.Header.Get("Authorization") != ""
}

// IsAuthorizationParameter checks if the request is using authorization parameter.
func IsAuthorizationParameter(r *http.Request) bool {
	return r.FormValue("access_token") != ""
}

// ExtractToken is the health check handler.
func ExtractToken(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		bearer = r.FormValue("access_token")
	}
	return strings.TrimPrefix(bearer, "Bearer ")
}
