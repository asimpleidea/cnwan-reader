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
	"fmt"
	"reflect"
	"testing"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCreateAWSConfigFromOpts(t *testing.T) {
	// TODO: this does not really work!

	a := assert.New(t)
	nop := zerolog.Nop()

	cases := []struct {
		opts option.CloudMap
		exp  []func(*config.LoadOptions) error
	}{
		{
			opts: option.CloudMap{},
			exp:  []func(*config.LoadOptions) error{},
		},
		// {
		// 	opts: option.CloudMap{
		// 		Region: "us-east-1",
		// 		Authentication: &option.CloudMapAuthentication{
		// 			Profile:             "test-profile",
		// 			CredentialsFilePath: "/some/path",
		// 		},
		// 	},
		// 	exp: []func(*config.LoadOptions) error{
		// 		config.WithRegion("us-east-1"),
		// 		config.WithSharedConfigProfile("test-profile"),
		// 		config.WithSharedConfigFiles([]string{"/some/path"}),
		// 	},
		// },
	}

	for i, currCase := range cases {
		cfg := createAWSConfigFromOpts(currCase.opts, nop, nop)

		for _, expC := range currCase.exp {
			rc := reflect.ValueOf(expC)
			found := false
			for _, retC := range cfg {
				if reflect.ValueOf(retC) == rc {
					found = true
				}
			}
			if !found {
				a.FailNow("a config is not found!", fmt.Sprintf("case %d", i))
			}
		}
	}
}
