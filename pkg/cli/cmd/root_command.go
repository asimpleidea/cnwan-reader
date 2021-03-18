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
	"os"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	log    zerolog.Logger         = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	vlevel map[bool]zerolog.Level = map[bool]zerolog.Level{false: zerolog.Disabled, true: zerolog.InfoLevel}
)

// NewRootCommand returns the root command.
func NewRootCommand() *cobra.Command {
	// -------------------------------
	// Define root command
	// -------------------------------

	globalFlags := option.ReaderOptions{}

	cmd := &cobra.Command{
		Use:   "cnwan-reader [command] [flags]",
		Short: "CN-WAN Reader observes changes in metadata in a service registry.",
		Long: `CN-WAN Reader connects to a service registry and
observes changes about registered services, delivering found events to a
a separate handler for processing.`,
		Example: "cnwan-reader poll cloudmap --interval 15",
	}

	// -------------------------------
	// Define global persistent flags
	// -------------------------------

	// TODO: sanitize localhost and make it docker.internal
	// TODO: redo descriptions
	cmd.PersistentFlags().BoolVar(&globalFlags.Verbose, "verbose", internal.DefaultVerbose, "if set, the logs will be more verbose")
	cmd.PersistentFlags().StringVar(&globalFlags.AdaptorURL, "adaptor-url", internal.DefaultAdaptorURL, "the url, in form of host:port/path, where the events will be sent to. Look at the documentation to learn more about this.")
	cmd.PersistentFlags().StringSliceVar(&globalFlags.RequiredMetadataKeys, "required-metadata-keys", []string{}, "a list of comma-separated metadata keys to look for. Example: traffic-profile,replicas")
	cmd.PersistentFlags().BoolVar(&globalFlags.IncludeOptionalKeys, "include-optional-keys", internal.DefaultIncludeOptionalKeys, "if set and the me")

	// -------------------------------
	// Define subcommands
	// -------------------------------

	cmd.AddCommand(newPollCommand(&globalFlags))

	return cmd
}
