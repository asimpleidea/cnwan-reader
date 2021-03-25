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
	"io"
	"os"
	"sync"

	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// logger       zerolog.Logger
	// vlogger      zerolog.Logger
	readerLogger *Logger
	once         sync.Once
)

type Logger struct {
	r zerolog.Logger
	v zerolog.Logger
}

func (l *Logger) Regular() zerolog.Logger {
	return l.r
}

func (l *Logger) Verbose() zerolog.Logger {
	return l.v
}

// GetLoggers returns two loggers:
// - a basic one that only logs important messages (from info above)
// - and another one used only for verbose logging (so only info)
//
// If opts.Verbose is false, v contains a nop logger that doesn't log anywhere.
// If opts.ConsoleOnly is true, no logs will be printed to file.
func GetLogger(opts option.Log) *Logger {
	once.Do(func() {
		logWriters := []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr}}
		if opts.LogsFilePath != "" {
			logWriters = append(logWriters,
				&lumberjack.Logger{
					Filename: opts.LogsFilePath,
					MaxSize:  DefaultLogsMaxSize,
					MaxAge:   DefaultLogsMaxDays,
				})
		}

		l := zerolog.New(io.MultiWriter(logWriters...)).With().Timestamp().Logger().Level(zerolog.InfoLevel)
		v := zerolog.Nop()
		if opts.Verbose {
			v = zerolog.New(io.MultiWriter(logWriters...)).With().Timestamp().Logger().Level(zerolog.InfoLevel)
		}

		readerLogger = &Logger{r: l, v: v}
	})

	return readerLogger
}

func FromSliceToMap(slice []string) (m map[string]bool) {
	m = map[string]bool{}
	for _, val := range slice {
		m[val] = true
	}

	return
}
