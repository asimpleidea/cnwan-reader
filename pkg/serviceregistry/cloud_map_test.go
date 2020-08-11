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
	"errors"
	"testing"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/openapi"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceData(t *testing.T) {
	a := &awsCloudMap{
		metadataKey: "key",
	}
	errNoKey := errors.New("instance does not have the required metadata key")
	errNoAddress := errors.New("instance has no address")
	errNoID := errors.New("instance has no id")

	// Case 1: no id
	inst := &servicediscovery.InstanceSummary{}
	res, err := a.getServiceData(inst)
	assert.Nil(t, res)
	assert.Equal(t, errNoID, err)

	var id string
	inst.Id = &id
	res, err = a.getServiceData(inst)
	assert.Nil(t, res)
	assert.Equal(t, errNoID, err)

	id = "id"
	inst.Id = &id

	// Case 2: Has no metadata
	res, err = a.getServiceData(inst)
	assert.Nil(t, res)
	assert.Equal(t, errNoKey, err)

	inst.Attributes = map[string]*string{}
	res, err = a.getServiceData(inst)
	assert.Nil(t, res)
	assert.Equal(t, errNoKey, err)

	var str string
	inst.Attributes[a.metadataKey] = &str
	res, err = a.getServiceData(inst)
	assert.Nil(t, res)
	assert.Equal(t, errNoKey, err)

	// Case 3: no addresses
	str = "val"
	inst.Attributes[a.metadataKey] = &str
	res, err = a.getServiceData(inst)
	assert.Nil(t, res)
	assert.Equal(t, errNoAddress, err)

	// Case 4: full data
	ip4 := "10.10.10.10"
	inst.Attributes[awsIPv4Attr] = &ip4

	servExp := &openapi.Service{
		Name:     id,
		Metadata: []openapi.Metadata{{Key: a.metadataKey, Value: str}},
		Address:  ip4,
		Port:     awsDefaultInstancePort,
	}

	res, err = a.getServiceData(inst)
	assert.Equal(t, servExp, res)

	delete(inst.Attributes, awsIPv4Attr)
	ip6 := "2001:db8:a0b:12f0::1"
	inst.Attributes[awsIPv6Attr] = &ip6
	servExp.Address = ip6

	port := "8080"
	inst.Attributes[awsPortAttr] = &port
	servExp.Port = 8080
	res, err = a.getServiceData(inst)
	assert.Equal(t, servExp, res)
}
