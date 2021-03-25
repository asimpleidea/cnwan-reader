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

package read

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/event"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/openapi"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/serviceregistry"
)

func Poll(opts option.Poll, readOpts option.Global) (err error) {
	// -------------------------------
	// Set ups
	// -------------------------------

	log := internal.GetLogger(readOpts.Log)
	v := log.Verbose().With().Logger()
	l := log.Regular().With().Logger()

	v.Info().Str("command", "Poll").Msg("starting...")
	mainCtx := context.Background()
	var sr serviceregistry.ServiceRegistry
	wg := sync.WaitGroup{}
	wg.Add(2)

	// -------------------------------
	// Set up service registry
	// -------------------------------

	// if opts.ServiceDirectory != nil {
	// 	v.Info().Msg("using Google Service Directory")
	// 	sdctx, canc := context.WithTimeout(mainCtx, time.Minute)
	// 	sr, err = serviceregistry.NewGoogleServiceDirectoryStateReader(sdctx, *opts.ServiceDirectory)
	// 	if err != nil {
	// 		// TODO: check if err is context
	// 		canc()
	// 		return err
	// 	}
	// 	canc()
	// 	v.Info().Msg("authenticated")
	// }
	// defer sr.CloseClient()

	// -------------------------------
	// Set up dispatcher
	// -------------------------------

	// ctx, canc := context.WithCancel(mainCtx)
	// dispatcher := event.NewDispatcher(event.DispatcherOptions{
	// 	Verbose:         readOpts.Verbose,
	// 	TimeOut:         time.Minute,
	// 	AdaptorEndpoint: readOpts.AdaptorURL,
	// })

	// go func() {
	// 	defer wg.Done()
	// 	dispatcher.Work(ctx)
	// 	l.Info().Msg("done worker")
	// }()

	evChan := make(chan event.Operation, 50)
	go func() {
		defer wg.Done()
		event.Work(evChan)
	}()

	// -------------------------------
	// Start polling
	// -------------------------------

	v.Info().Str("interval", opts.Interval.String()).Msg("starting poller...")
	ctx, canc := context.WithCancel(mainCtx)
	// TODO: check that this gets canceled with CTRL+C
	go func(o option.Poll) {
		defer wg.Done()
		// TODO: define way to to log poll x100
		v.Info().Str("timeout", time.Minute.String()).Msg("getting initial state...")

		firstStateCtx, firstStateCanc := context.WithTimeout(ctx, time.Minute)
		firstState, err := sr.GetCurrentState(firstStateCtx, []string{})
		if err != nil {
			if err == context.Canceled {
				l.Error().Str("timeout", time.Minute.String()).Msg("timeout expired while trying to get initial state")
			} else {
				l.Error().Err(err).Msg("error occurred while trying to get initial state")
			}

			firstStateCanc()
			return
		}
		firstStateCanc()

		for _, serv := range firstState.Services {
			evChan <- event.Operation{
				Method:  event.CreateOperation,
				Service: event.Service{},
			}
			_ = serv
		}

		ticker := time.NewTicker(o.Interval)

		for {
			select {
			case <-ticker.C:
				// go func() {
				v.Info().Msg("polling...")
				if _, err := sr.GetCurrentState(firstStateCtx, []string{}); err != nil {
					l.Err(err).Msg("error while reading current state; skipping...")
				}
				// }()
			case <-ctx.Done():
				ticker.Stop()
			}
		}

	}(opts)

	// -------------------------------
	// Graceful shutdown
	// -------------------------------

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	fmt.Println()

	// Cancel the context and wait for objects that use it to receive
	// the stop command
	canc()
	wg.Wait()

	l.Info().Msg("good bye!")
	return nil
}

var i = 0

func getCurrentState(ctx context.Context, sr serviceregistry.ServiceRegistry, readOpts option.Global, dispatcher *event.Dispatcher) (*serviceregistry.CurrentState, error) {
	// currentState, err := sr.GetCurrentState(ctx, readOpts.RequiredMetadataKeys)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, ep := range currentState.Endpoints {
	// 	sv := currentState.Services[ep.ServName]
	// 	dispatcher.Enqueue(&openapi.Event{
	// 		Event: "create",
	// 		Service: openapi.Service{
	// 			Name:    ep.Name,
	// 			Address: ep.Address,
	// 			Port:    ep.Port,
	// 			Metadata: func() (m []openapi.Metadata) {
	// 				for k, v := range sv.Metadata {
	// 					m = append(m, openapi.Metadata{Key: k, Value: v})
	// 				}
	// 				return
	// 			}(),
	// 		},
	// 	})
	// }
	fmt.Println("got current", i)
	dispatcher.Enqueue(&openapi.Event{
		Event: "create",
		Service: openapi.Service{
			Name:    fmt.Sprintf("%d", i),
			Address: fmt.Sprintf("%d", i),
			Port:    int32(i),
		},
	})
	i++
	fmt.Println("after state", i-1)

	return nil, nil
}
