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

// NewRootCommand defines the root command, including its flags definition
// and subcommands and returns it, so that it could be used.
func NewRootCommand() *cobra.Command {
	globalOpts := option.Global{}

	// -------------------------------
	// Define root command
	// -------------------------------

	cmd := &cobra.Command{
		Use: `cnwan-reader poll|watch SERVICE_REGISTRY|SERVICE [OPTIONS]
	[--help, -h]`,

		Short: `CN-WAN Reader observes a service registry and detects changes
to registered services.`,

		Long: `CN-WAN Reader connects to a service registry or a third party
provider and observes changes occurring to registered services.
Detected changes are then parsed and delivered to an external component for
processing called Adaptor.

Observation is performed with one of these two methods -- the one you choose
depends on which service registry/provider you want to observe:

	- polling (cnwan-reader poll ...), which consists of periodic API calls to
	  a service registry to get its current state and see what has changed.
	  To get more help about this method and see what is supported, please
	  execute "cnwan-reader poll --help".
	- watching (cnwan-reader watch ...), that uses the service
	  registry's/application's library to receive changes as they happen.
	  To get more help about this method and see what is supported, please
	  execute "cnwan-reader watch --help".

Detected changes are parsed and delivered to an Adaptor, a software that
receives the changes and performs some kind of processing with the data it
received. An example of an adaptor is the CN-WAN Adaptor, which sends the data
to an SD-WAN controller. Make sure you define the correct URL for the adaptor
with the --adaptor-url option.`,

		Example: "cnwan-reader poll cloudmap -k traffic-profile -i 15s",

		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if err := parseAdaptorURL(&globalOpts.AdaptorURL); err != nil {
				return err
			}

			if globalOpts.Log.LogsFilePath != "" {
				if err := isLogsFilePathValid(globalOpts.Log.LogsFilePath); err != nil {
					return err
				}
			}

			return nil
		},
	}

	// -------------------------------
	// Define global persistent flags
	// -------------------------------

	cmd.PersistentFlags().StringVarP(&globalOpts.AdaptorURL, "adaptor-url", "u", internal.DefaultAdaptorURL, "the url of the adaptor, where data will be sent to. Write as scheme://host:port/path .")
	cmd.PersistentFlags().BoolVarP(&globalOpts.Log.Verbose, "verbose", "v", internal.DefaultVerbose, "if set, the logs will be more verbose and contain more lines.")
	cmd.PersistentFlags().StringVarP(&globalOpts.Log.LogsFilePath, "logs-file", "l", "", "the path to the log file, including its name.")
	cmd.PersistentFlags().BoolVar(&globalOpts.DryRun, "dry-run", false, "if set, no data will be sent to the Adaptor, but will only be printed on console and/or log file.")

	// -------------------------------
	// Define subcommands
	// -------------------------------

	cmd.AddCommand(newPollCommand(&globalOpts))
	cmd.AddCommand(newWatchCommand(&globalOpts))
	cmd.AddCommand(newVersionCommand(&globalOpts))

	return cmd
}
