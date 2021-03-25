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
	"strings"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/spf13/cobra"
)

// newWatchEtcdCommand defines the etcd command and returns
// it so that it could be used as a subcommand.
func newWatchEtcdCommand(globalOpts *option.Global) *cobra.Command {
	etcdOpts := option.Etcd{}
	etcdAuth := option.EtcdAuthentication{}

	// -------------------------------
	// Define etcd command
	// -------------------------------

	cmd := &cobra.Command{

		Use: `etcd --metadata-keys,-k=KEY_1[,KEY_2,...]
	[--endpoints=ENDPOINT_1[,ENDPOINT_2,...]] [--user=USER] [--password=PASSWORD]
	[--prefix=PREFIX] [--help]`,

		Short: "watch changes to services posted to etcd.",

		Aliases: []string{"et"},

		Long: `etcd establishes a long watching connection to etcd and
detects changes as they happen, without the need for periodic checks.

In order to work and be able to detect services properly, your etcd should
contain CN-WAN's etcd service registry format: this means that you can of
course store everything that you want there, but you must plan to reserve
some keys only for CN-WAN's service registry. It is highly recommended to read
https://developer.cisco.com/docs/cloud-native-sdwan/#!etcd
and learn about these concepts before starting. If you are unsure, please
include --dry-run to see if things are working properly.

The only truly required option is --metadata-keys, which provides the filter to
apply when monitoring services. Only services that have *all* the metadata keys
you provide will be observed.

In order to successfully connect to etcd, a list of endpoints where etcd is
running is needed. While you don't have to enter **all** of them, make sure
to enter a few, so to be able to maintain the watching connection alive in
case one of them goes briefly down. The list of endpoints should be in the
form of host:port. 

If you added roles to etcd, you will need to provide the user with --user and
its password with --password, or the authentication will fail. Currently,
authenticating through certificates is not supported yet.

OPTIONS

	--metadata-keys provides the filter to apply when monitoring services.
	Only services that have *all* the annotations keys you provide here
	will be observed. This is the only **required** option.

	--endpoints is a comma-separated list of endpoints in the form of host:ip
	where your etcd nodes are serving from. For example, if you have only two
	machines running etcd, you may write IP_MACHINE_1:PORT,IP_MACHINE_2:PORT.
	Just as a reminder, etcd's default port is 2379, in case you didn't change
	it. Make sure to enter a few if you have multiple nodes, so to be able to
	keep the watching connection alive in case one of them goes briefly down.
		
	--prefix is an optional prefix to use and insert before each etcd key.
	By default, CN-WAN uses /service-registry and if you followed CN-WAN's
	documentation to create an etcd cluster then this is is already the way
	your etcd should be configured as, and you can leave this option blank.
	If for some reason you have another key, then please enter it here.
	Keep in mind that this will sanitized: multiple leading slashes (/) will be
	removed, as well as all trailing slashes (/). For example, if you enter
	"///my-key/", this value will be sanitized later as /my-key .
	
	--user defines the user to authenticate as. It's good practice to add an
	etcd user and assign a role to it. If you leave this empty CN-WAN Reader
	will try to authenticate as a guest. Please read the documentation on this
	repository to learn how to create a user and role for the CN-WAN Reader,
	giving it only the bare minimum amount of permissions it needs to operate.
	
	--password defines the password to use for the user defined with --user.
	Cannot be used without --user, and if a user is not provided then the
	execution will stop and prompt an error.`,

		Example: `etcd -k traffic-profile \
	--prefix /another/prefix/service-registry`,

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := persistentPreRunE(cmd, args); err != nil {
				return err
			}

			if len(etcdOpts.MetadataKeys) == 0 {
				return fmt.Errorf("no metadata keys provided")
			}

			if len(etcdOpts.Endpoints) == 0 {
				return fmt.Errorf("not endpoints provided")
			}

			if etcdOpts.Prefix != "" {
				etcdOpts.Prefix = strings.TrimPrefix(etcdOpts.Prefix, "//")
				etcdOpts.Prefix = strings.TrimSuffix(etcdOpts.Prefix, "/")
			}

			if etcdAuth.User == "" {
				if etcdAuth.Password == "" {
					return nil
				}

				return fmt.Errorf("password provided but --user not provided")
			}

			etcdOpts.Authentication = &etcdAuth
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement me
			// read.Watch(globalOpts, WithEtcd(etcdOpts))
			return nil
		},
	}

	// -------------------------------
	// Define flags
	// -------------------------------

	cmd.Flags().StringSliceVar(&etcdOpts.Endpoints, "endpoints", []string{"localhost:2379"}, "a list of endpoints where your etcd nodes are listening from.")
	cmd.Flags().StringVar(&etcdOpts.Prefix, "prefix", "/service-registry", "the key prefix where to watch keys from.")
	cmd.Flags().StringVar(&etcdAuth.User, "user", "", "the user to authenticate as.")
	cmd.Flags().StringVar(&etcdAuth.Password, "password", "", "the password to use for the user. Must be used with --user.")
	cmd.Flags().StringSliceVarP(&etcdOpts.MetadataKeys, "metadata-keys", "k", []string{}, "a list of comma-separated metadata keys to look for. Example: traffic-profile,replicas .")

	return cmd
}
