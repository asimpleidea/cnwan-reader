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
	"github.com/spf13/cobra"
)

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloudmapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloudmapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCloudMap(cmd *cobra.Command, args []string) {
	// TODO...
}
