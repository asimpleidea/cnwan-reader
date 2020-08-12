// Copyright © 2020 Cisco
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// All rights reserved.

/*
 * CNWAN Reader API
 *
 * The CNWAN Reader implements the [service discovery](https://en.wikipedia.org/wiki/Service_discovery) pattern by connecting to a service registry and observing changes in registered services/endpoints. Detected changes are then processed and sent as events to the API endpoints defined below.  Events are **sent** to the following endpoints, thus any program interested in receiving them must generate the *server* code from this OpenAPI specification and define their own logic in the generated code.  By default, the CNWAN Reader expects the server that will receive events to operate on port `80` and receive events on `/cnwan/events`, but if your server uses a different port/endpoint you can override this value on the generated server code with the one your server is using. Once done, when launching the CNWAN Reader specify the correct endpoint by providing it as a command line argument, e.g. with `--adaptor-api localhost:9909` events will be sent on `localhost:9909/events`, and with `--adaptor-api example.com/another/path` events will be sent to `example.com/another/path/events`.  As a final note, please take in mind that this specification can also serve as a reference/guide for the creation of an adaptor.   As a matter of fact, your adaptor can even provided its own OpenAPI which includes the endpoints described here with different descriptions and different meanings for the response codes, or it can even include other endpoints as well. But as long as formats, returned response code and the endpoints of this specification match the ones on your adaptor's specification, compatibility with CNWAN Reader is guaranteed.
 *
 * API version: 1.0.0 beta
 * Contact: cnwan@cisco.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"
)

// APIResponse stores the API response returned by the server.
type APIResponse struct {
	*http.Response `json:"-"`
	Message        string `json:"message,omitempty"`
	// Operation is the name of the OpenAPI operation.
	Operation string `json:"operation,omitempty"`
	// RequestURL is the request URL. This value is always available, even if the
	// embedded *http.Response is nil.
	RequestURL string `json:"url,omitempty"`
	// Method is the HTTP method used for the request.  This value is always
	// available, even if the embedded *http.Response is nil.
	Method string `json:"method,omitempty"`
	// Payload holds the contents of the response body (which may be nil or empty).
	// This is provided here as the raw response.Body() reader will have already
	// been drained.
	Payload []byte `json:"-"`
}

// NewAPIResponse returns a new APIResonse object.
func NewAPIResponse(r *http.Response) *APIResponse {

	response := &APIResponse{Response: r}
	return response
}

// NewAPIResponseWithError returns a new APIResponse object with the provided error message.
func NewAPIResponseWithError(errorMessage string) *APIResponse {

	response := &APIResponse{Message: errorMessage}
	return response
}
