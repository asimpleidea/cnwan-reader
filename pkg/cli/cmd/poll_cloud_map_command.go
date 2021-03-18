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
func newPollCloudMapCmd(pollOpts *option.PollOptions) *cobra.Command {
	localFlags := option.CloudMapOptions{
		Authentication: &option.CloudMapAuthenticationOptions{},
	}
	optionsPath := ""

	// -------------------------------
	// Unmarshal and validate functions
	// -------------------------------

	// unmarshal the file pointed by optionsPath and set the values found there
	// to localFlags, unless already set.
	// unmarshal := func() error {
	// 	options := option.CloudMapOptions{}

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

	// 	if len(localFlags.Region) == 0 {
	// 		localFlags.Region = options.Region
	// 	}

	// 	if len(localFlags.CredentialsPath) == 0 {
	// 		localFlags.CredentialsPath = options.CredentialsPath
	// 	}

	// 	return nil
	// }

	// // validate data in localFlags
	// validate := func() error {
	// 	if len(localFlags.Region) == 0 {
	// 		return fmt.Errorf("no region set")
	// 	}

	// 	return nil
	// }

	// -------------------------------
	// Define poll command
	// -------------------------------

	// TODO: expand long
	cmd := &cobra.Command{
		Use:   "cloudmap [flags]",
		Short: "connect to Cloud Map to get registered services",
		Long: `cloudmap connects to AWS CloudMap and
observes changes to services published in it, i.e. metadata, addresses and
ports.
	
For this to work, a valid region must be provided with --region and the
aws credentials must be properly set.

Unless a different credentials path is defined with --credentials-path,
$HOME/.aws/credentials on Linux/Unix and %USERPROFILE%\.aws\credentials on
Windows will be used instead. Alternatively, credentials path can be set
with environment variables. For a complete list of alternatives, please
refer to AWS Session documentation, but, to keep things simple, we suggest you
use the default one.
Run --help to get a description of all the flags.`,
		Example: "cnwan-reader poll cloudmap --region us-west-2 --credentials path/to/credentials/file",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			// -------------------------------
			// Parse the options
			// -------------------------------

			// if len(optionsPath) > 0 {
			// 	v.Info().Str("options-path", optionsPath).Msg("parsing options file...")
			// 	if err := unmarshal(); err != nil {
			// 		log.Err(err).Msg("could not parse options file; skipping...")
			// 	}
			// }

			// -------------------------------
			// Validate
			// -------------------------------

			if len(localFlags.Region) == 0 {
				return fmt.Errorf("no region set")
			}

			if len(localFlags.Authentication.CredentialsPath) == 0 {
				localFlags.Authentication = nil
			}
			// if err := validate(); err != nil {
			// 	log.Err(err).Msg("error while parsing options; exiting...")
			// 	return cmd.Help()
			// }

			pollOpts.CloudMap = &localFlags
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement me
			return nil
		},
	}

	// -------------------------------
	// Define flags
	// -------------------------------

	cmd.Flags().StringVar(&localFlags.Region, "region", "", "gcloud region location. Example: us-west-2")
	cmd.Flags().StringVar(&localFlags.Authentication.CredentialsPath, "credentials-path", "", "path to aws credentials. Example: ./credentials")
	cmd.Flags().StringVar(&optionsPath, "options", "", "the path to the yaml file containing options")

	return cmd
}
