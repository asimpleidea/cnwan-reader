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

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/spf13/cobra"
)

// TODO: maybe pass pollFlags instad of global
func newPollServiceDirectoryCmd(pollOpts *option.PollOptions) *cobra.Command {
	localFlags := option.ServiceDirectoryOptions{
		Authentication: &option.ServiceDirectoryAuthenticationOptions{},
	}
	optionsPath := ""

	// -------------------------------
	// Unmarshal and validation functions
	// -------------------------------

	// unmarshal the file pointed by optionsPath and set the values found there
	// to localFlags, unless already set.
	// unmarshal := func() error {
	// 	options := option.ServiceDirectoryOptions{}

	// 	f, err := os.Open(optionsPath)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer f.Close()

	// 	if err := yaml.NewDecoder(f).Decode(&options); err != nil {
	// 		return err
	// 	}

	// 	// -------------------------------
	// 	// Use data from options, if not set
	// 	// -------------------------------

	// 	if len(localFlags.ProjectID) == 0 {
	// 		localFlags.ProjectID = options.ProjectID
	// 	}

	// 	if len(localFlags.Region) == 0 {
	// 		localFlags.Region = options.Region
	// 	}

	// 	if len(localFlags.ServiceAccountPath) == 0 {
	// 		localFlags.ServiceAccountPath = options.ServiceAccountPath
	// 	}

	// 	if len(localFlags.APIKey) == 0 {
	// 		localFlags.APIKey = options.APIKey
	// 	}

	// 	return nil
	// }

	// -------------------------------
	// Define poll command
	// -------------------------------

	// TODO: expand long
	cmd := &cobra.Command{
		Use:   "servicedirectory [flags]",
		Short: "connect to Service Directory to get registered services",
		Long: `This command connects to Google Cloud Service Directory and
observes changes in services published in it, i.e. metadata, addresses and
ports.

In order to work, a project, location and valid credentials must be provided.
Run --help to get a description of all the flags.`,
		Example: "cnwan-reader poll servicedirectory --region us-west2 --api-key 12345ABCDEF",

		RunE: func(cmd *cobra.Command, args []string) error {
			// -------------------------------
			// Parse the options
			// -------------------------------

			// TODO: implement this part
			// if len(optionsPath) > 0 {
			// 	v.Info().Str("options-path", optionsPath).Msg("parsing options file...")
			// 	if err := unmarshal(); err != nil {
			// 		log.Err(err).Msg("could not parse options file; skipping...")
			// 	}
			// }

			// -------------------------------
			// Validate
			// -------------------------------

			if len(localFlags.ProjectID) == 0 {
				return fmt.Errorf("no project ID set")
			}

			if len(localFlags.Region) == 0 {
				return fmt.Errorf("no region set")
			}

			if len(localFlags.Authentication.APIKey) == 0 && len(localFlags.Authentication.ServiceAccountPath) == 0 {
				localFlags.Authentication = nil
			}

			pollOpts.ServiceDirectory = &localFlags
			return nil
		},
	}

	// -------------------------------
	// Define flags
	// -------------------------------

	cmd.Flags().StringVar(&localFlags.ProjectID, "project-id", "", "the project ID to use")
	cmd.Flags().StringVar(&localFlags.Region, "region", "", "gcloud region location. Example: us-west2")
	cmd.Flags().StringVar(&localFlags.Authentication.ServiceAccountPath, "service-account", "", "path to the gcp service account. Example: ./service-account.json")
	cmd.Flags().StringVar(&localFlags.Authentication.APIKey, "api-key", "", "api-key to use for authentication")
	cmd.Flags().StringVar(&optionsPath, "options", "", "the path to the yaml file containing options")

	return cmd
}

// func optionsFromFile(optionsPath string) (*option.ReaderOptions, error) {
// 	f, err := os.Open(optionsPath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	var options option.ReaderOptions
// 	if err := yaml.NewDecoder(f).Decode(&options); err != nil {
// 		return nil, err
// 	}

// 	return &options, nil
// }

// TODO: expand this
// func mergeReaderFlags(cmd *cobra.Command, cli, file *option.ReaderOptions) (*option.ReaderOptions, error) {
// 	if !cmd.Flag("verbose").Changed {
// 		cli.Verbose = file.Verbose
// 	}

// 	if !cmd.Flag("adaptor-url").Changed && len(file.AdaptorURL) > 0 {
// 		cli.AdaptorURL = file.AdaptorURL
// 	}

// 	if !cmd.Flag("required-metadata-keys").Changed && len(file.RequiredMetadataKeys) > 0 {
// 		cli.RequiredMetadataKeys = file.RequiredMetadataKeys
// 	}

// 	if !cmd.Flag("").Changed && len(file.RequiredMetadataKeys) > 0 {
// 		cli.RequiredMetadataKeys = file.RequiredMetadataKeys
// 	}

// 	return cli, nil
// }

// func sanitizeLocalhost(u string) (string, error) {

// }
