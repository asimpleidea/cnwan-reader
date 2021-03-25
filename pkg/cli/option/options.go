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

// Global contains options that are inherited by all commands and therefore
// are valid for any given command.
type Global struct {
	// AdaptorURL contains the URL that must be used to communicate with the
	// Adaptor.
	AdaptorURL string `yaml:"adaptorURL"`
	// DryRun specifies whether detected changes will be sent or not.
	DryRun bool `yaml:"dryRun,omitEmpty"`
	// Log contains log options.
	Log Log `yaml:"logOptions,omitempty"`
	// Poll contains options for the poll command
	Poll *Poll `yaml:"poll,omitempty"`
}

// Log contains logging options.
type Log struct {
	// Verbose specifies whether more log lines should be produced.
	Verbose bool `yaml:"verbose,omitempty"`
	// LogsFilePath is the path where the logs will be written.
	LogsFilePath string `yaml:"logsFilePath,omitempty"`
}

// Poll contains options about the poll command and that will be inherited
// by its sub-commands.
type Poll struct {
	// Interval defines how frequent the polling must be.
	Interval time.Duration `yaml:"interval,omitempty"`
}

// Watch contains options about the watch command and that will be inherited
// by its sub-commands.
type Watch struct {
	// As of now, there is no configuration for watch.
}

// CloudMap contains options to connect and get data from AWS Cloud Map.
type CloudMap struct {
	// Region where to get data from.
	Region string `yaml:"region"`
	// Authentication contains options for authenticating to AWS.
	Authentication *CloudMapAuthentication `yaml:",inline"`
	// AttributeKeys is a list of required keys to look for.
	AttributesKeys []string `yaml:"attributesKeys,omitempty"`
}

// CloudMapAuthentication contains data for authenticating to AWS.
type CloudMapAuthentication struct {
	// Profile to authenticate as.
	Profile string `yaml:"profile"`
	// CredentialsPath is the path to the aws credentials file to
	// use to authenticate.
	CredentialsFilePath string `yaml:"credentialsFile"`
}

// ServiceDirectory sets the options to connect to and get data
// from Google Service Directory.
type ServiceDirectory struct {
	// Region where to get data from.
	Region string `yaml:"region,omitempty"`
	// ProjectID is the ID of the google cloud project.
	ProjectID string `yaml:"projectID,omitempty"`
	// Authentication contains authentication modes for GCP
	Authentication *ServiceDirectoryAuthentication `yaml:",inline"`
	// AnnotationsKeys is an array of containing all the keys that are
	// required for a service to be monitored.
	AnnotationsKeys []string `yaml:"annotationsKeys"`
}

type ServiceDirectoryAuthentication struct {
	// ServiceAccountPath is the path to the google service account to
	// use to authenticate.
	ServiceAccountPath string `yaml:"serviceAccountPath,omitempty"`
}

// Etcd contains options for authenticating to etcd and retrieving
// information from it.
type Etcd struct {
	// Endpoints is a list of ip:port of the nodes where etcd is running.
	Endpoints []string `yaml:"endpoints"`
	// Prefix is an optional value that etcd keys will have. This will be
	// inserted before any monitored key.
	Prefix string `yaml:"prefix,omitempty"`
	// MetadataKeys is a list of keys that must be present inside a service's
	// metadata in order to be monitored.
	MetadataKeys []string `yaml:"metadataKeys"`
	// Authentication contains options for authenticating to etcd.
	Authentication *EtcdAuthentication `yaml:",inline"`
}

// EtcdAuthentication contains options for authenticating to etcd.
type EtcdAuthentication struct {
	// User to authenticate as.
	User string `yaml:"user,omitempty"`
	// Password to use for the user specified on the field above.
	Password string `yaml:"password,omitempty"`
}

// K8s contains options to connect to a Kubernetes cluster and get data
// from it.
type K8s struct {
	// AnnotationsKeys is a list of annotation keys to look for in a service.
	AnnotationsKeys []string `yaml:"annotationsKeys"`
	// Kubeconfig contains data about how to connect to the Kubernetes cluster,
	// for example the path to the kubeconfig file or the context to use.
	KubeConfig *KubeConfig `yaml:",inline"`
}

// KubeConfig contains options on how to connect to the Kubernetes cluster.
type KubeConfig struct {
	// Path to the kubeconfig file.
	Path string `yaml:"kubeConfigPath,omitempty"`
	// Context to use authenticatication.
	Context string `yaml:"context,omitempty"`
}
