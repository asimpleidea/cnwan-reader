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
	"context"
	"fmt"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/serviceregistry"
	"github.com/spf13/cobra"
)

// newPollCloudMapCmd defines the cloud map command and returns it so that it
// could be used as a subcommand.
func newPollCloudMapCmd(globalOpts *option.Global, pollOpts *option.Poll) *cobra.Command {
	cmOpts := option.CloudMap{}
	cmAuth := option.CloudMapAuthentication{}
	var metKeys []string

	// -------------------------------
	// Define cloudmap command
	// -------------------------------

	cmd := &cobra.Command{
		Use: `cloudmap --attribute-keys,-k=KEY_1[,KEY_2,...] 
	[--credentials-path=PATH] [--region,-r=REGION] [--profile,-p=PROFILE]
	[--help,-h]`,

		Short: "connect to AWS Cloud Map to get registered services.",

		Aliases: []string{"cm"},

		Long: `cloudmap connects to AWS CloudMap and observes changes to
registered services and their instances.

The only truly required option is --attribute-keys, or --metadata-keys which
is only an alias to it, to provide the filter to apply when monitoring
services. Only services that have *all* the attributes keys you provide will
be observed.

CN-WAN Reader needs to authenticate to CloudMap to work, and a credentials file
must be used to do that.
If you have the aws cli installed and want to use the default
configuration you can leave all options empty, apart from --attribute-keys,
as all the other options will just override configurations. The default
credentials path is $HOME/.aws/credentials on Linux/Unix and
%USERPROFILE%\.aws\credentials on Windows.

You can read
https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html
to learn how to set environment variables if your machines already have them
set or if you prefer to do that instead.

OPTIONS:

	--attributes-keys provides the filter to apply when monitoring services.
		Only services that have *all* the attributes keys you provide here
		will be observed. This is the only **required** option.

	--metadata-keys is just an alias for --attributse-keys.

	--credentials-file-path specifies the credentials file to be used for
		authenticating to AWS. Please consult AWS documentation on credentials
		files to understand how they are created and how they work.
		Fill this option if you want to use a credentials file different from
		the default one, e.g. the ones created by the aws cli, if installed.

	--profile is used to override the default profile and authenticate as
		another profile. If not set, the default one inside the provided
		credentials file will be used, i.e. the one used by the aws cli,
		if installed.

	--region overrides the region to use. If not set, the default region from
		the provided credentials file will be used. Use this in case the
		CloudMap services are on a region different from the one in your
		provided credentials file.`,

		Example: `cloudmap --k traffic-profile --region us-west-2 \
	--credentials path/to/credentials/file`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := persistentPreRunE(cmd, args); err != nil {
				return err
			}

			if len(cmOpts.AttributesKeys) == 0 {
				if len(metKeys) == 0 {
					return fmt.Errorf("no attribute keys provided")
				}
				cmOpts.AttributesKeys = metKeys
			}

			if cmAuth.CredentialsFilePath != "" || cmAuth.Profile != "" {
				cmOpts.Authentication = &cmAuth
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: finish implementation for this.
			a, err := serviceregistry.NewAWSCloudMapStateReader(context.Background(), cmOpts, globalOpts.Log)
			if err != nil {
				return err
			}

			a.GetCurrentState(context.Background())
			// return Poll(globalOpts, pollOpts, option.WithCloudMap())
			return nil
		},
	}

	// -------------------------------
	// Define flags
	// -------------------------------

	cmd.Flags().StringVar(&cmAuth.CredentialsFilePath, "credentials-file-path", "", "path to the AWS credentials.")
	cmd.Flags().StringVarP(&cmOpts.Region, "region", "r", "", "AWS region location. Example: us-west-2")
	cmd.Flags().StringVar(&cmAuth.Profile, "profile", "", "the AWS profile to use.")
	cmd.Flags().StringSliceVarP(&cmOpts.AttributesKeys, "attributes-keys", "k", []string{}, "a comma-separated list of attributes keys to look for in services, e.g.: traffic-profile,stage.")
	cmd.Flags().StringSliceVar(&metKeys, "metadata-keys", []string{}, "alias for --attribute-keys. If --attribute-keys is present, this is ignored.")

	return cmd
}
