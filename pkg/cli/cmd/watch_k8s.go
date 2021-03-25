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
	"github.com/spf13/cobra"
)

// newWatchK8sCommand defines the kubernetes command and returns
// it so that it could be used as a subcommand.
func newWatchK8sCommand(globalOpts *option.Global) *cobra.Command {
	kOpts := option.K8s{}
	kubeconf := option.KubeConfig{}

	// -------------------------------
	// Define kubernetes command
	// -------------------------------

	cmd := &cobra.Command{

		Use: `kubernetes --annotations-keys,-k=KEY_1[,KEY_2,...]
	[--kubeconfig=KUBECONFIG_PATH] [--context=CONTEXT] [--help]`,

		Short: "observe changes to services deployed to a Kubernetes cluster",

		Aliases: []string{"k8s"},

		Long: `kubernetes command establishes a long watching connection to
a Kubernetes cluster and detectes changes to LoadBalancer type of services
as they happend, without requiring periodic checks.

To establish a connection, a valid kubeconfig file must be provided: if you
already havea  kubeconfig file on your computer, e.g. in your machine's default
directory, you can use that one and its default context without specifying its
path or the context. You may override those values with --kubeconfig and/or
--context.

The only truly required option is --annotations-keys, which provides the filter
to apply when monitoring services. Only services that have *all* the annotations
keys you provide will be observed.

For more information, it is higly recommended to learn more about these
concepts on CN-WAN's documentation files.

OPTIONS

	--annotations-keys provides the filter to apply when monitoring services.
	Only services that have *all* the annotations keys you provide here
	will be observed. This is the only **required** option.
	
	--kubeconfig is the path to the kubeconfig file. Usually this is retrieved
	automatically, e.g. if you have used kubectl successfully before or you
	created/added your cluster with your provider's cli tool. If you don't want
	to use the default one or it was not possible to retrieve it, you may fill
	this with the path to your desired kubeconfig file.
	
	--context is the context name to use on the provided --kubeconfig path or
	the default kubeconfig. If not set, the default context for the found
	kubeconfig will be used. If you want to watch another cluster from that
	kubeconfig, you may specify the context here.`,

		Example: "kubernetes --annotations-keys traffic-profile --context gke",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := persistentPreRunE(cmd, args); err != nil {
				return err
			}

			if kubeconf.Context != "" || kubeconf.Path != "" {
				kOpts.KubeConfig = &kubeconf
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	// -------------------------------
	// Define flags
	// -------------------------------

	cmd.Flags().StringVar(&kubeconf.Context, "context", "", "the kubeconfig context to use.")
	cmd.Flags().StringVar(&kubeconf.Path, "kubeconfig", "", "path to the kubeconfig file to use.")
	cmd.Flags().StringSliceVarP(&kOpts.AnnotationsKeys, "annotations-keys", "k", []string{}, "a list of comma-separated annotation keys to look for. Example: traffic-profile,replicas")

	return cmd
}
