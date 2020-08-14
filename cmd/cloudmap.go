// Copyright © 2020 Cisco
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
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var awsRegion string

// cloudmapCmd represents the cloudmap command
var cloudmapCmd = &cobra.Command{
	Use:   "cloudmap",
	Short: "Connect to Cloud Map to get registered services",
	Long: `This command connects to AAWS Cloud Map and
	observes changes in services published in it, i.e. metadata, addresses and
	ports.
	
	In order to work, location and valid credentials must be provided.`,
	Run:     runCloudMap,
	Aliases: []string{"cm", "aws"},
}

func init() {
	rootCmd.AddCommand(cloudmapCmd)

	cloudmapCmd.Flags().StringVar(&awsRegion, "region", "", "aws region location. Example: us-west-2")
}

func runCloudMap(cmd *cobra.Command, args []string) {
	var err error
	l := log.With().Str("func", "cmd.runCloudMap").Logger()
	l.Info().Msg("starting...")

	ctx, canc := context.WithCancel(context.Background())

	// Parse flags
	if len(awsRegion) == 0 {
		l.Fatal().Err(fmt.Errorf("%s", "region not provided")).Msg("fatal error encountered")
	}
	if len(credsPath) > 0 {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", credsPath)
	}

	// TODO...
}
