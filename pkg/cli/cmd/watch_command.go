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
	_ "github.com/CloudNativeSDWAN/cnwan-reader/pkg/read"
	"github.com/spf13/cobra"
)

// newWatchCommand defines the watch command, including its flags and
// subcommands.
func newWatchCommand(globalOpts *option.Global) *cobra.Command {
	// -------------------------------
	// Define watch command
	// -------------------------------

	cmd := &cobra.Command{
		Use: "watch etcd|kubernetes [OPTIONS] [--help]",

		Short: "observe changes with 'watch' methods.",

		Aliases: []string{"wa"},

		Long: `watch observes changes to a service registry or provider
by utilizing methods that are not based on periodic checks but rather use
methods defined by the service registry's or provider's library/sdk.

Each time a change has been made, the change is detected -- almost immediately
in most cases -- without having to periodically check the entire service
registry.

Currently, only etcd and Kubernetes are supported, with cnwan-reader watch etcd
and cnwan-reader watch kubernetes respectively. To get more information about
either of those, include --help when you run them.`,

		Example: `watch etcd -k traffic-profile \
	--endpoints 10.10.10.10:2379,10.10.10.11:2379`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := persistentPreRunE(cmd, args); err != nil {
				return err
			}

			return nil
		},
	}

	// -------------------------------
	// Define sub commands flags
	// -------------------------------

	cmd.AddCommand(newWatchEtcdCommand(globalOpts))
	cmd.AddCommand(newWatchK8sCommand(globalOpts))

	return cmd
}
