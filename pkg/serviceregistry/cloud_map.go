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
	"errors"
	"fmt"
	"strconv"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/openapi"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/rs/zerolog/log"
)

const (
	awsIPv4Attr            string = "AWS_INSTANCE_IPV4"
	awsIPv6Attr            string = "AWS_INSTANCE_IPV6"
	awsPortAttr            string = "AWS_INSTANCE_PORT"
	awsDefaultInstancePort int32  = 80
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
	l := log.With().Str("func", "Handler.GetServices").Logger()
	maps := map[string]*openapi.Service{}

	// First, get services
	services, err := a.getServicesIDs()
	if err != nil {
		l.Err(err).Msg("error while getting services")
	}

	for _, servID := range services {
		l = l.With().Str("service-id", servID).Logger()

		// Get the instances
		instances, err := a.getInstances(servID)
		if err != nil {
			l.Err(err).Msg("error while getting instances, skipping...")
			continue
		}

		for _, inst := range instances {
			data, err := a.getServiceData(inst)
			if err != nil {
				l.Debug().Msg(err.Error())
				continue
			}

			if data != nil {
				l.Debug().Str("instance-id", data.Name).Msg("instance has the required metadata key")
				mapKey := fmt.Sprintf("%s_%d", data.Address, data.Port)
				maps[mapKey] = data
			}
		}
	}

	return maps
}

func (a *awsCloudMap) getServicesIDs() ([]string, error) {
	l := log.With().Str("func", "Handler.awsCloudMap.getServicesIDs").Logger()

	out, err := a.sd.ListServices(&servicediscovery.ListServicesInput{})
	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, errors.New("received nil response")
	}

	servIDs := []string{}
	for _, service := range out.Services {
		if service.Id == nil || (service.Id != nil && len(*service.Id) == 0) {
			l.Debug().Msg("a service with no/empty ID has been found and is going to be skipping...")
			continue
		}

		servIDs = append(servIDs, *service.Id)
	}

	return servIDs, nil
}

func (a *awsCloudMap) getInstances(servID string) ([]*servicediscovery.InstanceSummary, error) {
	out, err := a.sd.ListInstances(&servicediscovery.ListInstancesInput{
		ServiceId: &servID,
	})
	if err != nil {
		return nil, err
	}

	if out == nil {
		return nil, errors.New("received nil response")
	}

	return out.Instances, nil
}

func (a *awsCloudMap) getServiceData(inst *servicediscovery.InstanceSummary) (*openapi.Service, error) {
	if inst.Id == nil || (inst.Id != nil && len(*inst.Id) == 0) {
		return nil, errors.New("instance has no id")
	}

	metadataVal := inst.Attributes[a.metadataKey]
	if metadataVal == nil || (metadataVal != nil && len(*metadataVal) == 0) {
		return nil, errors.New("instance does not have the required metadata key")
	}

	address := ""
	if ipv6 := inst.Attributes[awsIPv6Attr]; ipv6 != nil && len(*ipv6) > 0 {
		address = *ipv6
	}
	if ipv4 := inst.Attributes[awsIPv4Attr]; ipv4 != nil && len(*ipv4) > 0 {
		address = *ipv4
	}
	if len(address) == 0 {
		return nil, errors.New("instance has no address")
	}

	var port int32 = awsDefaultInstancePort
	if instancePort := inst.Attributes[awsPortAttr]; instancePort != nil {
		intPort, err := strconv.ParseInt(*instancePort, 10, 32)
		if err == nil {
			port = int32(intPort)
		}
	}

	return &openapi.Service{
		Address:  address,
		Name:     *inst.Id,
		Metadata: []openapi.Metadata{{Key: a.metadataKey, Value: *metadataVal}},
		Port:     port,
	}, nil
}
