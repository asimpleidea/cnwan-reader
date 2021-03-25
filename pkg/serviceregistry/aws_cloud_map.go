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

package serviceregistry

import (
	"context"
	"errors"
	"fmt"
	"time"

	optr "github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/smithy-go"
	"github.com/rs/zerolog"
)

type StateReader struct {
}

type StateRecorder struct {
	Service string
}

type AWSCloudMapStateReader struct {
	log  zerolog.Logger
	vlog zerolog.Logger
	keys map[string]bool
	opts optr.CloudMap
	cfg  aws.Config
}

func NewAWSCloudMapStateReader(ctx context.Context, opts optr.CloudMap, lopts optr.Log) (*AWSCloudMapStateReader, error) {
	// -------------------------------
	// Init and setups
	// -------------------------------

	l, v := func() (zerolog.Logger, zerolog.Logger) {
		log := internal.GetLogger(lopts)
		return log.Regular().With().Str("name", "CloudMapStateReader").Logger(), log.Verbose().With().Str("name", "CloudMapStateReader").Logger()
	}()

	// -------------------------------
	// Load configuration
	// -------------------------------

	v.Info().Msg("loading AWS credentials...")
	cfg, err := config.LoadDefaultConfig(ctx, createAWSConfigFromOpts(opts, l, v)...)
	if err != nil {
		var apierr smithy.APIError
		if errors.As(err, &apierr) {
			return nil, fmt.Errorf("timeout expired while trying load AWS credentials")
		}

		return nil, err
	}

	// Note: after some tests, I discovered that config.LoadDefaultConfig
	// **never** fails, even if you run this on an ec2 machine with neither
	// default aws files nor ec2 roles set. For now, let's just check if the
	// region is found and if not, let's return an error.
	_, err = cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	if cfg.Region == "" {
		return nil, fmt.Errorf("could not retrieve any AWS region")
	}

	l.Info().Msg("AWS credentials successfully retrieved")

	return &AWSCloudMapStateReader{
		log:  l,
		vlog: v,
		keys: internal.FromSliceToMap(opts.AttributesKeys),
		opts: opts,
		cfg:  cfg,
	}, nil
}

func createAWSConfigFromOpts(opts optr.CloudMap, log, vlog zerolog.Logger) (configs []func(*config.LoadOptions) error) {
	configs = []func(*config.LoadOptions) error{}

	if opts.Region != "" {
		vlog.Info().Str("region", opts.Region).Msg("using region defined on command options")
		configs = append(configs, config.WithRegion(opts.Region))
	}

	if opts.Authentication == nil {
		return configs
	}

	if opts.Authentication.Profile != "" {
		vlog.Info().Str("profile", opts.Authentication.Profile).Msg("using profile defined on command options")
		configs = append(configs, config.WithSharedConfigProfile(opts.Authentication.Profile))
	}

	if opts.Authentication.CredentialsFilePath != "" {
		vlog.Info().Msg("using credentials file defined on command options")
		configs = append(configs, config.WithSharedConfigFiles([]string{opts.Authentication.CredentialsFilePath}))
	}

	return
}

// TODO: pass channel for sending what is found
func (a *AWSCloudMapStateReader) GetCurrentState(ctx context.Context) {
	listCtx, canc := context.WithTimeout(ctx, time.Minute)
	servs, err := a.GetServices(listCtx)
	if err != nil {
		canc()
		return
	}
	canc()

	for _, serv := range servs {

		go func(servID string) {
			instCtx, cancel := context.WithTimeout(ctx, time.Minute)
			// TODO: this needs to be changed
			instances, err := a.GetInstances(instCtx, aws.ToString(&servID))
			if err != nil {
				a.log.Info().Err(err).Str("service-id", servID).Msg("could not get instances for this service, skipping...")
				cancel()
				return
			}
			cancel()

			for _, inst := range instances {
				// TODO: put to channel
				a.log.Info().Str("inst-id", inst).Msg("got instance")
			}
		}(serv)

	}
}

// TODO: this should actually return a list of Service, not slice of strings
func (a *AWSCloudMapStateReader) GetServices(ctx context.Context) ([]string, error) {
	client := servicediscovery.NewFromConfig(a.cfg)
	input := servicediscovery.ListServicesInput{}

	out, err := client.ListServices(ctx, &input)
	if err != nil {
		return nil, err
	}

	a.vlog.Debug().Int("#", len(out.Services)).Msg("got services")

	ids := make([]string, len(out.Services))
	for i, serv := range out.Services {
		ids[i] = aws.ToString(serv.Id)
	}

	return ids, nil
}

// TODO: change this to serviceSummary
func (a *AWSCloudMapStateReader) GetInstances(ctx context.Context, serviceID string) ([]string, error) {
	client := servicediscovery.NewFromConfig(a.cfg)
	input := servicediscovery.ListInstancesInput{
		ServiceId: aws.String(serviceID),
	}

	out, err := client.ListInstances(ctx, &input)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(out.Instances))
	for i, inst := range out.Instances {
		ids[i] = aws.ToString(inst.Id)
	}

	return ids, nil
}

func (a *AWSCloudMapStateReader) Close() {
	// No need to close the client
}
