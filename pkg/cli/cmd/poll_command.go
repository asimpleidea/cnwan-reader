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
	"time"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	_ "github.com/CloudNativeSDWAN/cnwan-reader/pkg/read"
	"github.com/spf13/cobra"
)

// newPollCommand defines the poll command, including its flags and the
// subcommands.
func newPollCommand(globalOpts *option.Global) *cobra.Command {
	pollOpts := option.Poll{}

	// -------------------------------
	// Define poll command
	// -------------------------------

	cmd := &cobra.Command{
		Use: "poll servicedirectory|cloudmap [OPTIONS] [--help]",

		Short: "periodically check a service registry for changes.",

		Aliases: []string{"po"},

		Long: `poll uses a polling mechanism to detect changes to a
service registry, that is performing periodic calls to the the service
registry's API to retrieve the current state of the services that are
registered in it.

This command must be used when wanting to check service registries that
do not provide their own watching/observation mechanisms (yet), such as
Google Service Directory or AWS Cloud Map.

The current sitation is then confronted with the previous state that was
cached inside the CN-WAN Reader and, if different, the new data is then sent
to the adaptor. For example, if a service is not found anymore, a DELETE is
sent to the adaptor.

Currently, only Google Service Directory and AWS Cloud Map are supported with
this command and must be used as cnwan-reader poll servicedirectory or
cnwan-reader poll cloudmap respectively. Whichever you intend to use, run
--help to get more information about it.

OPTIONS:

	--poll-interval must be used to set the duration between two consecutive
	checks: it accepts durations in a human-friendly manner, such as 1m for
	1 minute, 30s for 30 seconds or even combinations such 1m20s if you want
	to perform checks every minute and 20 seconds. What value you choose
	depends on how frequent changes to the services are made and/or how quickly
	you want to react to them. You cannot set values lower than 1 second (1s).`,

		Example: `poll cloudmap --k traffic-profile --region us-west-2 \
	--credentials path/to/credentials/file`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := persistentPreRunE(cmd, args); err != nil {
				return err
			}

			if pollOpts.Interval < time.Second {
				return fmt.Errorf("provided polling interval is not valid")
			}

			return nil
		},
	}

	// -------------------------------
	// Define persistent flags
	// -------------------------------

	cmd.PersistentFlags().DurationVarP(&pollOpts.Interval, "poll-interval", "i", internal.DefaultPollInterval, "interval between two consecutive requests.")

	// -------------------------------
	// Define sub commands flags
	// -------------------------------

	cmd.AddCommand(newPollServiceDirectoryCmd(globalOpts, &pollOpts))
	cmd.AddCommand(newPollCloudMapCmd(globalOpts, &pollOpts))

	return cmd
}
