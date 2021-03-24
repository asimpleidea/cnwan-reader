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
	"path"
	"strings"
	"time"

	sd "cloud.google.com/go/servicedirectory/apiv1beta1"
	op "github.com/CloudNativeSDWAN/cnwan-operator/pkg/servregistry"
	optr "github.com/CloudNativeSDWAN/cnwan-reader/pkg/cli/option"
	"google.golang.org/api/iterator"
	optg "google.golang.org/api/option"
	sdpb "google.golang.org/genproto/googleapis/cloud/servicedirectory/v1beta1"
)

type CurrentState struct {
	Namespaces map[string]*op.Namespace // namespace-name => namespace
	Services   map[string]*op.Service   // service-name => service
	Endpoints  map[string]*op.Endpoint  // endpoint-name => endpoint
}

type GoogleServiceDirectory struct {
	client   *sd.RegistrationClient
	basePath string
}

func NewGoogleServiceDirectoryReader(ctx context.Context, opts optr.ServiceDirectoryOptions) (*GoogleServiceDirectory, error) {
	// TODO: this version does not implement APIKey and will be introduced on
	// next PR/release.
	// TODO: this version still uses v1beta1 and will use v1 on next PR/release.
	// TODO: on next version perform a TestIAMPermissions before doing anything else?
	var canc context.CancelFunc
	if ctx == nil {
		ctx, canc = context.WithTimeout(context.Background(), time.Minute)
		defer canc()
	}

	cli, err := sd.NewRegistrationClient(ctx, generateOptsForGSD(opts)...)
	if err != nil {
		return nil, err
	}

	return &GoogleServiceDirectory{
		client:   cli,
		basePath: path.Join("projects", opts.ProjectID, "locations", opts.Region),
	}, nil
}

func (g *GoogleServiceDirectory) GetCurrentState(ctx context.Context, keys []string) (*CurrentState, error) {
	// -------------------------------
	// Init
	// -------------------------------

	keyFilter := func() string {
		if len(keys) == 0 {
			return ""
		}

		filters := make([]string, len(keys))
		for i := 0; i < len(keys); i++ {
			filters[i] = fmt.Sprintf("metadata.%s!=''", keys[i])
		}

		return strings.Join(filters, " AND ")
	}()

	var canc context.CancelFunc
	if ctx == nil {
		ctx, canc = context.WithTimeout(context.Background(), time.Minute)
		defer canc()
	}

	namespaces := map[string]*op.Namespace{}
	services := map[string]*op.Service{}
	endpoints := map[string]*op.Endpoint{}

	// -------------------------------
	// Get namespaces
	// -------------------------------

	req := &sdpb.ListNamespacesRequest{Parent: g.basePath}
	it := g.client.ListNamespaces(ctx, req)
	if it == nil {
		return nil, fmt.Errorf("no namespaces")
	}

	for {
		nsResp, err := it.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}

			fmt.Println("err ns list", err)
			continue
		}

		// -------------------------------
		// Get services
		// -------------------------------

		srvReq := &sdpb.ListServicesRequest{Parent: nsResp.Name, Filter: keyFilter}
		sv := g.client.ListServices(ctx, srvReq)
		if sv == nil {
			continue
		}

		for {
			svResp, err := sv.Next()
			if err != nil {
				if err == iterator.Done {
					break
				}

				break
			}

			// -------------------------------
			// Get endpoints
			// -------------------------------

			epReq := &sdpb.ListEndpointsRequest{Parent: svResp.Name}
			ep := g.client.ListEndpoints(ctx, epReq)

			for {
				epResp, err := ep.Next()
				if err != nil {
					if err == iterator.Done {
						break
					}

					break
				}

				endpoints[epResp.Name] = &op.Endpoint{
					Name:     epResp.Name,
					NsName:   nsResp.Name,
					ServName: svResp.Name,
					Address:  epResp.Address,
					Port:     epResp.Port,
					Metadata: func() map[string]string {
						// TODO: copy other metadata?
						if len(epResp.Metadata) > 0 {
							return epResp.Metadata
						}

						return map[string]string{}
					}(),
				}
			}

			services[svResp.Name] = &op.Service{
				Name:   svResp.Name,
				NsName: nsResp.Name,
				Metadata: func() map[string]string {
					// TODO: copy other metadata?
					if len(svResp.Metadata) > 0 {
						return svResp.Metadata
					}

					return map[string]string{}
				}(),
			}

			// Insert the namespace. Here to avoid inserting it if the
			// namespace had no services.
			if _, exists := namespaces[nsResp.Name]; !exists {
				namespaces[nsResp.Name] = &op.Namespace{
					Name: nsResp.Name,
					Metadata: func() map[string]string {
						// TODO: copy other metadata?
						if len(nsResp.Labels) > 0 {
							return nsResp.Labels
						}

						return map[string]string{}
					}(),
				}
			}
		}
	}

	return &CurrentState{
		Endpoints:  endpoints,
		Services:   services,
		Namespaces: namespaces,
	}, nil
}

func (g *GoogleServiceDirectory) CloseClient() {
	if g.client != nil {
		g.client.Close()
	}
}

func generateOptsForGSD(sdOpts optr.ServiceDirectoryOptions) (opts []optg.ClientOption) {
	opts = []optg.ClientOption{}
	if sdOpts.Authentication == nil {
		return
	}

	auth := sdOpts.Authentication
	if len(auth.ServiceAccountPath) > 0 {
		opts = append(opts, optg.WithCredentialsFile(auth.ServiceAccountPath))
	}

	if len(auth.APIKey) > 0 {
		opts = append(opts, optg.WithAPIKey(auth.APIKey))
	}

	return
}
