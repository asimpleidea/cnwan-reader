// Copyright © 2021 Cisco
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

package event

import "time"

// DispatcherOptions contains settings to fine tune the behavior of the
// event dispatcher.
type DispatcherOptions struct {
	// Verbose specifies whether to display more log lines.
	Verbose bool
	// TimeOut defines the maximum time to wait for a response from the
	// CN-WAN Adaptor before terminating the call, to avoid unresponsiveness.
	TimeOut time.Duration
	// AdaptorEndpoint is the address where to send events to.
	AdaptorEndpoint string

	// TODO: add other options
}

type OperationMethod string

const (
	CreateOperation OperationMethod = "CREATE"
	UpdateOperation OperationMethod = "UPDATE"
	DeleteOperation OperationMethod = "DELETE"
)

type Operation struct {
	Method  OperationMethod
	Service Service
}

type Service struct {
}
