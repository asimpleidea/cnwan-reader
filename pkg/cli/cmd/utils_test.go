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

package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	"github.com/stretchr/testify/assert"
)

func TestParseAdaptorURL(t *testing.T) {
	a := assert.New(t)

	cases := []struct {
		url    string
		mode   string
		expURL string
		expErr error
	}{
		{
			url:    "myaddress.com",
			expURL: "http://myaddress.com",
		},
		{
			url:    "aaa/test",
			expURL: "",
			expErr: func() error {
				_, err := url.Parse("aaa")
				return err
			}(),
		},
		{
			url:    "aaa/test",
			expURL: "",
			expErr: func() error {
				_, err := url.Parse("aaa")
				return err
			}(),
		},
		{
			url:    "localhost:80/cnwan/",
			expURL: "http://localhost:80/cnwan",
		},
		{
			url:    "localhost:80/cnwan/",
			mode:   "docker",
			expURL: "http://host.docker.internal:80/cnwan",
		},
		{
			url:    "127.0.0.1:8080/cnwan/",
			mode:   "docker",
			expURL: "http://host.docker.internal:8080/cnwan",
		},
	}

	failed := func(i int) {
		a.FailNow("case failed", fmt.Sprintf("case %d", i))
	}
	for i, currCase := range cases {
		os.Setenv("MODE", currCase.mode)
		err := parseAdaptorURL(&currCase.url)
		if !a.Equal(currCase.expURL, currCase.expURL) || !a.Equal(currCase.expErr, err) {
			failed(i)
		}
	}
}

func TestIsLogsFilePathValid(t *testing.T) {
	a := assert.New(t)
	anyErr := fmt.Errorf("any")

	cases := []struct {
		prerun  func()
		postrun func()
		path    string
		expErr  error
	}{
		{
			path:   "/",
			expErr: anyErr,
		},
		{
			path: func() string {
				hdir := internal.DefaultHomeDirectory()
				return path.Join(hdir, ".cnwan-temp", "reader", "logs")
			}(),
			postrun: func() {
				hdir := internal.DefaultHomeDirectory()
				os.RemoveAll(path.Join(hdir, ".cnwan-temp"))
			},
		},
		{
			path:   "/text",
			expErr: anyErr,
		},
		{
			path:   "/some/kind/of/test",
			expErr: anyErr,
		},
		{
			prerun: func() {
				os.Create("testing")
			},
			path: "testing",
			postrun: func() {
				os.Remove("testing")
			},
		},
		{
			path: "testing",
			postrun: func() {
				os.Remove("testing")
			},
		},
	}

	for i, currCase := range cases {
		if currCase.prerun != nil {
			currCase.prerun()
		}
		err := isLogsFilePathValid(currCase.path)
		if err == nil {
			if currCase.expErr != nil {
				a.FailNow("was expecting error but none occurred", fmt.Sprintf("case %d", i))
			}
		} else {
			if currCase.expErr == nil {
				a.FailNow("wasn't expecting error but did occur", fmt.Sprintf("case %d", i))
			}
		}
		if currCase.postrun != nil {
			currCase.postrun()
		}
	}
}
