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
	"fmt"
	"os"
	"path"

	"cloud.google.com/go/compute/metadata"
	servicedirectory "cloud.google.com/go/servicedirectory/apiv1"
	op "github.com/CloudNativeSDWAN/cnwan-operator/pkg/servregistry"
	optr "github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"github.com/CloudNativeSDWAN/cnwan-reader/pkg/internal"
	"github.com/rs/zerolog"
	aug "golang.org/x/oauth2/google"
	optg "google.golang.org/api/option"
	servicedirectorypb "google.golang.org/genproto/googleapis/cloud/servicedirectory/v1"
)

// TODO: this may not be useful
type CurrentState struct {
	Namespaces map[string]*op.Namespace // namespace-name => namespace
	Services   map[string]*op.Service   // service-name => service
	Endpoints  map[string]*op.Endpoint  // endpoint-name => endpoint
}

type GoogleServiceDirectoryStateReader struct {
	client *servicedirectory.RegistrationClient
	opts   optr.ServiceDirectory
	log    zerolog.Logger
	vlog   zerolog.Logger
	keys   map[string]bool
}

func NewGoogleServiceDirectoryStateReader(ctx context.Context, opts optr.ServiceDirectory, lopts optr.Log) (*GoogleServiceDirectoryStateReader, error) {
	// -------------------------------
	// Init and setups
	// -------------------------------

	l, v := func() (zerolog.Logger, zerolog.Logger) {
		log := internal.GetLogger(lopts)
		return log.Regular().With().Str("name", "ServiceDirectoryStateReader").Logger(), log.Verbose().With().Str("name", "ServiceDirectoryStateReader").Logger()
	}()

	// -------------------------------
	// Get the client
	// -------------------------------

	v.Info().Msg("loading credentials...")
	cli, err := servicedirectory.NewRegistrationClient(ctx, createSDOpts(opts, l, v)...)
	if err != nil {
		return nil, err
	}
	l.Info().Msg("successfully retrieved credentials")

	// -------------------------------
	// Build the base path
	// -------------------------------

	if opts.ProjectID != "" {
		v.Info().Str("project-id", opts.ProjectID).Msg("using project ID provided with options")
	} else {
		if opts.Authentication != nil && opts.Authentication.ServiceAccountPath != "" {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", opts.Authentication.ServiceAccountPath)
		}

		v.Info().Msg("no project ID set, trying to retrieve it implictly...")
		cred, err := aug.FindDefaultCredentials(context.Background())
		if err != nil {
			return nil, err
		}
		opts.ProjectID = cred.ProjectID
		v.Info().Str("project-id", cred.ProjectID).Msg("successfully retrieved projectID")
	}

	if opts.Region != "" {
		v.Info().Str("region", opts.Region).Msg("using region provided with options")
	} else {
		v.Info().Msg("no region set, trying to retrieve it implicitly...")

		if !metadata.OnGCE() {
			return nil, fmt.Errorf("could not retrieve region")
		}

		v.Info().Msg("retrieving default region for project...")
		defRegion, err := metadata.ProjectAttributeValue("google-compute-default-region")
		if err != nil {
			return nil, fmt.Errorf("could not get default region from project")
		}
		v.Info().Str("google-compute-default-region", defRegion).Msg("successfully retrieved region")
		opts.Region = defRegion
	}

	return &GoogleServiceDirectoryStateReader{
		client: cli,
		opts:   opts,
		log:    l,
		vlog:   v,
		keys:   internal.FromSliceToMap(opts.AnnotationsKeys),
	}, nil
}

func (g *GoogleServiceDirectoryStateReader) buildPath(res ...string) (resPath string) {
	resPath = path.Join("projects", g.opts.ProjectID, "locations", g.opts.Region)

	if len(res) > 0 {
		resPath = path.Join(resPath, "namespaces", res[0])

		if len(res) > 1 {
			resPath = path.Join(resPath, "services", res[1])

			if len(res) > 2 {
				resPath = path.Join(resPath, "endpoints", res[2])
			}
		}
	}

	return
}

func (g *GoogleServiceDirectoryStateReader) GetCurrentState(ctx context.Context) {
	ns := g.client.ListNamespaces(context.Background(), &servicedirectorypb.ListNamespacesRequest{
		Parent: g.buildPath(),
	})
	for {
		nextNs, err := ns.Next()
		if err != nil {
			g.log.Info().Msg(err.Error())
			return
		}

		g.log.Info().Msg(nextNs.Name)
	}

	// TODO: implement me
}

func (g *GoogleServiceDirectoryStateReader) Close() {
	g.client.Close()
}

func createSDOpts(opts optr.ServiceDirectory, l, v zerolog.Logger) (gopts []optg.ClientOption) {
	gopts = []optg.ClientOption{}

	if opts.Authentication == nil {
		v.Info().Msg("performing implicit authentication...")
		return
	}

	if opts.Authentication.ServiceAccountPath != "" {
		v.Info().Str("service-account-path", opts.Authentication.ServiceAccountPath).Msg("using service account provided from options")
		gopts = []optg.ClientOption{optg.WithCredentialsFile(opts.Authentication.ServiceAccountPath)}
		return
	}

	return
}
