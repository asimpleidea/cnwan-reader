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

package servicedirectory

import (
	"fmt"
	"os"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/configuration"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var (
	log zerolog.Logger
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	log = zerolog.New(output).With().Timestamp().Logger().Level(zerolog.InfoLevel)
}

// GetServiceDirectory returns the servicedirectory command
//
// TODO: on next version this will probably be changed and adopt some
// other programming pattern, maybe with a factory.
func GetServiceDirectoryCommand() *cobra.Command {
	var sd *gcloudServiceDirectory
	options := &options{}
	var confPath string

	cmd := &cobra.Command{
		Use:   "servicedirectory",
		Short: "Connect to Service Directory to get registered services",
		Long: `This command connects to Google Cloud Service Directory and
		observes changes in services published in it, i.e. metadata, addresses and
		ports.

		In order to work, a project, location and valid credentials must be provided.`,
		Example: "servicedirectory --region us-west2 --service-account-path path/to/service_account.json",
		PreRun: func(cmd *cobra.Command, _ []string) {
			opts, err := parseFlags(cmd, configuration.GetConfigFile())
			if err != nil {
				log.Fatal().Err(err).Msg("fatal error encountered")
				return
			}

			if len(opts.credsPath) > 0 {
				os.Setenv("AWS_SHARED_CREDENTIALS_FILE", opts.credsPath)
			}

			if opts.debug {
				log = log.Level(zerolog.DebugLevel)
			}

			sess, err := session.NewSession()
			if err != nil {
				log.Fatal().Err(err).Msg("could not start AWS session")
				return
			}
			sd := servicediscovery.New(sess, aws.NewConfig().WithRegion(opts.region))

			cm = &awsCloudMap{
				opts: opts,
				sd:   sd,
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			run(cm)
		},
	}

	// Flags
	cmd.Flags().StringVar(&confPath, "conf", "", "path to the file containing configuration for Service Directory")
	cmd.Flags().StringVar(&options.region, "region", "", "region to use")
	cmd.Flags().StringVar(&options.projectID, "project-id", "", "the gcloud project id")
	cmd.Flags().StringVar(&options.serviceAccountPath, "service-account-path", "", "the path to the credentials file")
	cmd.Flags().StringSliceVar(&options.keys, "metadata-keys", []string{}, "the metadata keys to watch for")

	return cmd
}

func parseFlags(cmd *cobra.Command, conf *configuration.Config) (*options, error) {
	opts := &options{}

	if conf == nil {
		conf = &configuration.Config{
			ServiceRegistry: &configuration.ServiceRegistrySettings{
				AWSCloudMap: &configuration.CloudMapConfig{},
			},
		}
	}
	cmConf := conf.ServiceRegistry.AWSCloudMap

	awsRegion, _ := cmd.Flags().GetString("region")
	if len(awsRegion) == 0 {
		if len(cmConf.Region) == 0 {
			return nil, fmt.Errorf("region not provided")
		}

		awsRegion = cmConf.Region
	}
	opts.region = awsRegion

	credsPath, _ := cmd.Flags().GetString("credentials-path")
	if len(credsPath) == 0 {
		if len(cmConf.CredentialsPath) > 0 {
			credsPath = cmConf.CredentialsPath
		}
	}
	opts.credsPath = credsPath

	pollInterval := 5
	if cmd.Flags().Changed("poll-interval") {
		_pollInterval, _ := cmd.Flags().GetInt("poll-interval")
		if _pollInterval > 0 {
			pollInterval = _pollInterval
		}
	} else {
		if cmConf.PollInterval > 0 {
			pollInterval = cmConf.PollInterval
		}
	}
	opts.interval = pollInterval

	keys, err := utils.GetMetadataKeysFromCmdFlags(cmd)
	if err != nil {
		return nil, err
	}
	opts.keys = keys

	adaptor, err := utils.GetAdaptorEndpointFromFlags(cmd)
	if err != nil {
		return nil, err
	}
	opts.adaptor = adaptor
	opts.debug = utils.GetDebugModeFromFlags(cmd)

	return opts, nil
}

// func run(cm *awsCloudMap) {
// 	log.Info().Str("service-registry", "Cloud Map").Str("adaptor", cm.opts.adaptor).Msg("starting...")

// 	ctx, canc := context.WithCancel(context.Background())

// 	datastore := services.NewDatastore()
// 	servsHandler, err := services.NewHandler(ctx, cm.opts.adaptor)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("error while trying to connect to aws cloud map")
// 	}
// 	sendQueue := queue.New(ctx, servsHandler)

// 	go func() {
// 		log.Info().Msg("getting initial state...")
// 		oaSrvs, err := cm.getCurrentState(ctx)
// 		if err != nil {
// 			log.Fatal().Err(err).Msg("error while getting initial state of cloud map")
// 			return
// 		}

// 		log.Info().Msg("done")
// 		if filtered := datastore.GetEvents(oaSrvs); len(filtered) > 0 {
// 			go sendQueue.Enqueue(filtered)
// 		}

// 		// Get the poller
// 		log.Info().Msg("observing changes...")
// 		poll := poller.New(ctx, cm.opts.interval)
// 		poll.SetPollFunction(func() {
// 			oaSrvs, err := cm.getCurrentState(ctx)
// 			if err != nil {
// 				log.Err(err).Msg("error while polling, skipping...")
// 				return
// 			}

// 			if filtered := datastore.GetEvents(oaSrvs); len(filtered) > 0 {
// 				log.Info().Msg("changes detected")
// 				go sendQueue.Enqueue(filtered)
// 			}
// 		})

// 		poll.Start()
// 	}()

// 	// Graceful shutdown
// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, os.Interrupt)

// 	<-sig
// 	fmt.Println()
// 	log.Info().Msg("exit requested")

// 	// Cancel the context and wait for objects that use it to receive
// 	// the stop command
// 	canc()

// 	log.Info().Msg("good bye!")
// }
