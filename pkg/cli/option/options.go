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

package option

import "time"

// TODO: maybe remove Options, as that is repetitive (the package is already called option)
// ReaderOptions contains options that are valid throughout the entire program.
// TODO: expand documentation here.
type ReaderOptions struct {
	Verbose              bool         `yaml:"verbose,omitempty"`
	AdaptorURL           string       `yaml:"adaptorURL"`
	RequiredMetadataKeys []string     `yaml:"requiredMetadataKeys,omitempty"`
	IncludeOptionalKeys  bool         `yaml:"includeOptionalKeys,omitempty"`
	Poll                 *PollOptions `yaml:"poll,omitempty"`
}

type PollOptions struct {
	Interval         time.Duration            `yaml:"interval,omitempty"`
	ServiceDirectory *ServiceDirectoryOptions `yaml:"googleServiceDirectory,omitempty"`
	CloudMap         *CloudMapOptions         `yaml:"awsCloudMap,omitempty"`
}

// ServiceDirectoryOptions sets the options to connect to and get data
// from Google Service Directory.
// Only one between APIKey and ServiceAccountPath should be used to
// authenticate to google cloud platform. If both are set, then
// only APIKey will be used.
type ServiceDirectoryOptions struct {
	// Region where to get data from.
	Region string `yaml:"region"`
	// ProjectID is the ID of the google cloud project.
	ProjectID      string                                 `yaml:"projectID"`
	Authentication *ServiceDirectoryAuthenticationOptions `yaml:"authentication,omitempty"`
}

type ServiceDirectoryAuthenticationOptions struct {
	// APIKey is the API Key to use to authenticate to google cloud.
	APIKey string `yaml:"apiKey,omitempty"`
	// ServiceAccountPath is the path to the google service account to
	// use to authenticate.
	ServiceAccountPath string `yaml:"serviceAccountPath,omitempty"`
}

// CloudMapOptions sets the options to connect to and get data from
// AWS Cloud Map.
type CloudMapOptions struct {
	// Region where to get data from.
	Region         string                         `yaml:"region"`
	Authentication *CloudMapAuthenticationOptions `yaml:"authentication,omitempty"`
}

type CloudMapAuthenticationOptions struct {
	// CredentialsPath is the path to the aws credentials file to
	// use to authenticate.
	CredentialsPath string `yaml:"credentialsPath"`
}
