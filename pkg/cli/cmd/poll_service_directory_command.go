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

// newPollServiceDirectoryCmd defines the service directory command and returns
// it so that it could be used as a subcommand.
func newPollServiceDirectoryCmd(globalOpts *option.Global, pollOpts *option.Poll) *cobra.Command {
	sdOpts := option.ServiceDirectory{}
	sdAuth := option.ServiceDirectoryAuthentication{}
	var metKeys []string

	// -------------------------------
	// Define servicedirectory command
	// -------------------------------

	cmd := &cobra.Command{
		Use: `servicedirectory --annotations-keys,-k=KEY_1[,KEY_2,..]
	[--project-id,-p=PROJECT_ID] [--region,-r=REGION]
	[--service-account=SERVICE_ACCOUNT_PATH] [--help]`,

		Short: "connect to Google Service Directory to get registered services.",

		Aliases: []string{"sd"},

		Long: `servicedirectory connects to Google Cloud Service Directory and
observes changes to published services and their endpoints.

The only truly required option is --annotations-keys, or --metadata-keys which
is only an alias to it, to provide the filter to apply when monitoring
services. Only services that have *all* the annotations keys you provide will
be observed.

CN-WAN Reader needs to authenticate to Service Directory to work, together with
a project ID and a region where Service Directory is enabled.
You can provide the path to the service account JSON file with
--service-account-path, a project ID with --project-id and a region with
--region.
If you don't provide these values, some attempts to retrieve them implictly
will be made: i.e. searching inside gcloud cli default folder for the default
credentials JSON file, environment variables or at the credentials injected
inside the machine if you are running CN-WAN Reader inside a virtual machine
in Google Cloud. You may override such values by providing either of
--service-account-path, --project-id and/or --region, but if they are empty and
they were not found implicitly then CN-WAN Reader will stop execution because
it can't go on without those values.

You can read https://cloud.google.com/iam/docs/service-accounts for more
information about service acounts and authenticating to Google Cloud.

OPTIONS:

	--annotations-keys provides the filter to apply when monitoring services.
	Only services that have *all* the annotations keys you provide here
	will be observed. This is the only **required** option.

	--metadata-keys is just an alias for --annotations-keys.

	--service-account-path specifies the service account JSON file to be used
	for authenticating to Google Cloud. Please consult Google Cloud
	documentation on service accounts to understand how they are created and
	how they work.
	Fill this option if you want to use a service account different from
	the default one, e.g. the ones created by the gcloud cli, if installed.

	--project-id is the ID -- **not** the name -- of the project where Service
	Directory is enabled. If empty, CN-WAN Reader will try to retrieve this
	implictly as specified above. Execution will stop if the provided project
	ID was not found or it was not possible to retrieve it implictly.

	--region is the Service Directory region where to monitor services on. If
	empty, a default region will be retrieved, i.e. the one where the virtual
	machine is hosted in -- provided that CN-WAN Reader is running inside a
	virtual machine in Google Cloud. Execution will fail if the region was not
	found or it was not possible to retrieve this implictly.`,

		Example: "servicedirectory -k traffic-profile --region us-west2",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := persistentPreRunE(cmd, args); err != nil {
				return err
			}

			if len(sdOpts.AnnotationsKeys) == 0 {
				if len(metKeys) == 0 {
					return fmt.Errorf("no annotation keys provided")
				}

				sdOpts.AnnotationsKeys = metKeys
			}

			if sdAuth.ServiceAccountPath != "" {
				sdOpts.Authentication = &sdAuth
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			a, err := serviceregistry.NewGoogleServiceDirectoryStateReader(context.Background(), sdOpts, globalOpts.Log)
			if err != nil {
				return err
			}
			a.GetCurrentState(context.Background())
			// read.Poll(gloabalOpts, logOpts, WithServiceDirectory(opts.ServiceDirectory))
			return nil
		},
	}

	// -------------------------------
	// Define flags
	// -------------------------------

	cmd.Flags().StringVarP(&sdOpts.ProjectID, "project-id", "p", "", "the project ID to use")
	cmd.Flags().StringVarP(&sdOpts.Region, "region", "r", "", "gcloud region location. Example: us-west2")
	cmd.Flags().StringVar(&sdAuth.ServiceAccountPath, "service-account-path", "", "path to the gcp service account. Example: ./service-account.json")
	cmd.Flags().StringSliceVarP(&sdOpts.AnnotationsKeys, "annotations-keys", "k", []string{}, "a list of comma-separated metadata keys to look for. Example: traffic-profile,replicas")
	cmd.Flags().StringSliceVar(&metKeys, "metadata-keys", []string{}, "an alias for --annotation-keys")

	return cmd
}
