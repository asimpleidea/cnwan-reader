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
	"os/signal"
	"strings"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/poller"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/queue"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/serviceregistry"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/services"
	"github.com/rs/zerolog"
	l "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	mainCtx     context.Context
	cancCtx     context.CancelFunc
	debugMode   bool
	interval    int
	srHandler   serviceregistry.Handler
	sendQueue   queue.Queue
	datastore   services.Datastore
	endpoint    string
	metadataKey string
	credsPath   string
	poll        poller.Poller
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cnwan-reader",
	Short: "CNWAN Reader observes changes in metadata in a service registry.",
	Long: `CNWAN Reader connects to a service registry and 
observes changes about registered services, delivering found events to a
a separate handler for processing.`,

	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		mainCtx, cancCtx = context.WithCancel(context.Background())
		initObserveData()
	},

	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		poll.Start()
		gracefulShutdown()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "whether to log debug lines")
	rootCmd.PersistentFlags().IntVarP(&interval, "interval", "i", 5, "number of seconds between two consecutive polls")
	rootCmd.PersistentFlags().StringVar(&metadataKey, "metadata-key", "profile", "name of the metadata key to look for")
	rootCmd.PersistentFlags().StringVar(&credsPath, "credentials", "", "path to the credentials file")
	rootCmd.PersistentFlags().StringVar(&endpoint, "adaptor-api", "localhost/cnwan", "the api, in forrm of host:port/path, where the events will be sent to. Look at the documentation to learn more about this.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// -- Configure logger
	l.Logger = l.Output(zerolog.ConsoleWriter{
		Out: os.Stdout,
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s\t|", i)
		},
	}).Level(func() zerolog.Level {
		if debugMode {
			return zerolog.DebugLevel
		}

		return zerolog.InfoLevel
	}())
}

func initObserveData() {
	// Get the datastore
	datastore = services.NewDatastore()

	// Sanitize the endpoint
	endpoint = sanitizeAdaptorEndpoint(endpoint)

	// Get the queue
	servsHandler, err := services.NewHandler(mainCtx, endpoint)
	if err != nil {
		l.Fatal().Err(err).Msg("error while trying to initialize a service handler")
	}
	sendQueue = queue.New(mainCtx, servsHandler)

	// Get the poller
	poll = poller.New(mainCtx, interval)
	poll.SetPollFunction(processData)
}

func sanitizeAdaptorEndpoint(endp string) string {
	endp = strings.Trim(endp, "/")

	if strings.HasPrefix(endp, "localhost") {
		// Replace localhost in case we are running insde docker
		if mode := os.Getenv("MODE"); len(mode) > 0 && mode == "docker" {
			return strings.Replace(endp, "localhost", "host.docker.internal", 1)
		}
	}

	return endp
}

func gracefulShutdown() {
	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	fmt.Println()

	// Cancel the context and wait for objects that use it to receive
	// the stop command
	cancCtx()

	l.Info().Msg("good bye!")
}

func processData() {
	data := srHandler.GetServices()

	events := datastore.GetEvents(data)
	if len(events) > 0 {
		sendQueue.Enqueue(events)
	}
}
