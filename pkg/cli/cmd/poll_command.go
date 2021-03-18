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
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	"github.com/spf13/cobra"
)

func newPollCommand(globalOpts *option.ReaderOptions) *cobra.Command {
	// -------------------------------
	// Define poll command
	// -------------------------------

	localFlags := option.PollOptions{}

	// TODO: expand long
	cmd := &cobra.Command{
		Use:   "poll [servicedirectory|cloudmap] [flags]",
		Short: "poll a service registry to discover changes",
		Long: `poll uses a polling mechanism to detect changes to a
service registry. This means that the CN-WAN Reader will perform http calls to
the service registry, parse the result and see the difference.

This method is implemented only for those service registries that do not
provide a better way to do this: include --help or -h to know which ones are
included under this command`,
		Example: "poll cloudmap --region us-west-2 --credentials path/to/credentials/file",
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			globalOpts.Poll = &localFlags
			// TODO: implement me
			return nil
		},
	}

	// -------------------------------
	// Define persistent flags
	// -------------------------------

	cmd.PersistentFlags().DurationVar(&localFlags.Interval, "poll-interval", internal.DefaultPollInterval, "interval between two consecutive polls")

	// -------------------------------
	// Define sub commands flags
	// -------------------------------

	cmd.AddCommand(newPollServiceDirectoryCmd(&localFlags))
	cmd.AddCommand(newPollCloudMapCmd(&localFlags))

	return cmd
}
