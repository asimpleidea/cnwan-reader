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

package serviceregistry

import (
	"context"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/openapi"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

type awsCloudMap struct {
	session     *session.Session
	sd          *servicediscovery.ServiceDiscovery
	metadataKey string
	region      string
}

// NewCloudMapHandler returns a handler for Cloud Map
func NewCloudMapHandler(ctx context.Context, region, metadataKey string) (Handler, error) {
	// Create a Session with a custom region
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	// Create the service discovery
	sd := servicediscovery.New(sess, aws.NewConfig().WithRegion(region))

	return &awsCloudMap{
		session:     sess,
		region:      region,
		sd:          sd,
		metadataKey: metadataKey,
	}, nil
}

func (a *awsCloudMap) GetServices() map[string]*openapi.Service {
	// TODO: Implement me
	return nil
}
