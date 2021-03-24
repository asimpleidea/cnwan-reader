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

package serviceregistry

import (
	"context"

	sd "cloud.google.com/go/servicedirectory/apiv1beta1"
	optr "github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
)

type AWSCloudMap struct {
	client *sd.RegistrationClient
}

func NewAWSCloudMapReader(ctx context.Context, opts optr.ServiceDirectoryOptions) (*AWSCloudMap, error) {
	// TODO: implement me
	return nil, nil
}

func (a *AWSCloudMap) GetCurrentState(ctx context.Context, keys []string) {
	// TODO: implement me
}

func (a *AWSCloudMap) CloseClient() {
	// TODO: implement me
}
