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

package internal

import (
	"os"
	"path"
	"time"

	"github.com/rs/zerolog/log"
)

// TODO: check these out
const (
	DefaultVerbose             bool          = false
	DefaultAdaptorURL          string        = "http://localhost:80/cnwan"
	DefaultPollInterval        time.Duration = 5 * time.Second
	DefaultPollTimeout         time.Duration = time.Minute
	DefaultIncludeOptionalKeys bool          = false
	DefaultLogsFileName        string        = "logs"
	DefaultLogsMaxSize         int           = 100
	DefaultLogsMaxDays         int           = 30
	DefaultAWSDirName          string        = ".aws"
	DefaultAWSCredsFileName    string        = "credentials"
	// DefaultAWSProfile          string        = "default"

)

func DefaultHomeDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warn().Err(err).Msg("could not get home directory, using current directory")
		home = "."
	}

	return path.Join(home, ".cnwan", "reader")
}
