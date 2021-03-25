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
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	if parent := cmd.Parent(); parent != nil && parent.PersistentPreRunE != nil {
		return parent.PersistentPreRunE(parent, args)
	}

	return nil
}

func parseAdaptorURL(adurl *string) error {
	parsedURL := strings.Trim(*adurl, "/")

	if !strings.HasPrefix(parsedURL, "https://") && !strings.HasPrefix(parsedURL, "http://") {
		parsedURL = fmt.Sprintf("http://%s", parsedURL)
	}

	parsed, err := url.Parse(parsedURL)
	if err != nil {
		adurl = nil
		return err
	}

	if mode := os.Getenv("MODE"); strings.ToLower(mode) == "docker" && (parsed.Hostname() == "localhost" || parsed.Hostname() == "127.0.0.1") {
		parsedURL = strings.Replace(parsed.String(), "localhost", "host.docker.internal", 1)
	}

	*adurl = parsedURL
	return nil
}

func isLogsFilePathValid(logspath string) error {
	stat, err := os.Stat(logspath)
	if err == nil {
		if stat.IsDir() {
			return fmt.Errorf("%s is a directory", logspath)
		}

		return nil
	}

	parentDir := path.Dir(logspath)
	if parentDir != "." && parentDir != "/" {
		if err := recursiveCreateDirectory(parentDir); err != nil {
			return err
		}
	}

	f, err := os.Create(logspath)
	if err == nil {
		f.Close()
		return nil
	}

	return err
}

func recursiveCreateDirectory(dir string) error {
	parent := path.Dir(dir)
	if parent != "." && parent != "/" {
		if err := recursiveCreateDirectory(parent); err != nil {
			return err
		}
	}

	err := os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}
