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

package event

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/openapi"
	"github.com/rs/zerolog"
)

// EventDispatcher dispatches events and delivers them to the CN-WAN Adaptor.
// It acts as a queue that holds data before sending it: if the dispatcher
// is already busy sending events, then each new event is enqueued and sent
// when the worker is free again.
// If the queue is empty, the data is sent immediately.
type EventDispatcher interface {
	// Enqueue puts the provided event in the dispatcher's queue.
	// The worker will send this event together with the other ones in the
	// queue as soon as it finishes processing the current ones.
	Put(*openapi.Event)
	// Work spins up a new dispatcher worker that loops on the dispatcher's
	// internal queue, blocking when it is empty, thus without wasting CPU
	// cycles. Whenever a new item is inserted in the queue, the worker is
	// woken up and tries to send the data right away.
	Work(context.Context)
}

// Dispatcher dispatches events and delivers them to the CN-WAN Adaptor.
// It acts as a queue that holds data before sending it: if the dispatcher
// is already busy sending events, then each new event is enqueued and sent
// when the worker is free again.
// If the queue is empty, the data is sent immediately.
type Dispatcher struct {
	log    zerolog.Logger
	vlog   zerolog.Logger
	queue  map[string]*openapi.Event
	wakeup chan bool
	lock   sync.Mutex
}

// NewDispatcher returns a new instance of an event dispatcher.
func NewDispatcher(opts DispatcherOptions) *Dispatcher {
	// -------------------------------
	// Parse options
	// -------------------------------

	if opts.TimeOut == 0*time.Second {
		opts.TimeOut = time.Minute
	}

	if len(opts.AdaptorEndpoint) == 0 {
		opts.AdaptorEndpoint = openapi.NewConfiguration().BasePath
	}

	// -------------------------------
	// Set up logs
	// -------------------------------

	output := zerolog.ConsoleWriter{
		Out: os.Stdout,
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("Dispatcher: %s", i)
		},
	}

	log := zerolog.New(output).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	vlog := zerolog.Nop()
	if opts.Verbose {
		vlog = log.With().Logger().Level(zerolog.InfoLevel)
	}

	vlog.Info().Str("dispatching-timeout", opts.TimeOut.String()).Str("adaptor-endpoint", opts.AdaptorEndpoint).Msg("options parsed")

	// -------------------------------
	// Return the dispatcher
	// -------------------------------

	return &Dispatcher{
		log:    log,
		vlog:   vlog,
		queue:  map[string]*openapi.Event{},
		wakeup: make(chan bool),
		lock:   sync.Mutex{},
		// TODO: set openapi...
	}
}

// Enqueue puts the provided event in the dispatcher's queue.
// The worker will send this event together with the other ones in the queue
// as soon as it finishes processing the current ones.
func (c *Dispatcher) Enqueue(newEvent *openapi.Event) {
	c.enqueueEvent(newEvent)
	c.wakeup <- true
}

func (c *Dispatcher) enqueueEvent(newEvent *openapi.Event) {
	l := c.log.With().Str("event", fmt.Sprintf("%+v", newEvent)).Logger()
	v := l.With().Logger()
	l.Info().Msg("enqueuing event")

	c.lock.Lock()
	defer c.lock.Unlock()

	existing, there := c.queue[newEvent.Service.Name]

	// -------------------------------
	// Insert if no/same prior event
	// -------------------------------

	if !there {
		c.queue[newEvent.Service.Name] = newEvent
		v.Info().Msg("event inserted")
		return
	}

	if newEvent.Event == existing.Event {
		c.queue[newEvent.Service.Name] = newEvent
		v.Info().Msg("replaced existing event with same type on queue")
		return
	}

	// -------------------------------
	// Amend prior events
	// -------------------------------

	if newEvent.Event == "delete" {
		if existing.Event == "create" {
			delete(c.queue, newEvent.Service.Name)
			l.Info().Msg("removed existing 'create' event from queue")
		} else {
			c.queue[newEvent.Service.Name] = newEvent
			v.Info().Msg("replaced existing event on queue")
		}

		return
	}

	if existing.Event == "delete" {
		c.queue[newEvent.Service.Name] = &openapi.Event{
			Event:   "update",
			Service: newEvent.Service,
		}
	} else {
		c.queue[newEvent.Service.Name].Service = newEvent.Service
	}
	v.Info().Msg("replaced existing event")
}

// Work spins up a new dispatcher worker that loops on the dispatcher's
// internal queue, blocking when it is empty, thus without wasting CPU
// cycles. Whenever a new item is inserted in the queue, the worker is
// woken up and tries to send the data right away.
// Note that this need to be started in another goroutine and that only
// one instance of work must run at the same time.
func (c *Dispatcher) Work(ctx context.Context) {
	l := c.log.With().Logger()
	l.Info().Msg("worker started")

	for ctx.Err() == nil {
		c.lock.Lock()

		// -------------------------------
		// Is queue empty?
		// -------------------------------
		l.Info().Msg("acquired lock")
		if len(c.queue) == 0 {

			// Sleep until someone wakes us up.
			// When it happens go to next iteration an acquire the lock again.
			c.lock.Unlock()
			l.Info().Msg("release lock and sleeping")
			select {
			case <-c.wakeup:
				continue
			case <-ctx.Done():
				continue
			}
		}

		l.Info().Int("load", len(c.queue)).Msg("dispatching events...")
		events := []*openapi.Event{}
		for _, event := range c.queue {
			events = append(events, event)
		}
		c.queue = map[string]*openapi.Event{}
		c.lock.Unlock()

		// TODO: send data
		time.Sleep(10 * time.Second)
		l.Info().Msg("finished")
	}

	l.Info().Msg("worker exited")
}

func Work(evChan chan Operation) {
	// TODO: do something with the event
	for ev := range evChan {
		_ = ev
	}
}
