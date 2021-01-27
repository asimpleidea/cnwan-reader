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

package etcd

import (
	"context"

	"go.etcd.io/etcd/clientv3"
)

type fakeKV struct {
	_put     func(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	_get     func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	_delete  func(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error)
	_compact func(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error)
	_do      func(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error)
	_txn     func(ctx context.Context) clientv3.Txn
}

func (f *fakeKV) Put(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return f._put(ctx, key, val, opts...)
}

func (f *fakeKV) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return f._get(ctx, key, opts...)
}

func (f *fakeKV) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return f._delete(ctx, key, opts...)
}

func (f *fakeKV) Compact(ctx context.Context, rev int64, opts ...clientv3.CompactOption) (*clientv3.CompactResponse, error) {
	return f._compact(ctx, rev, opts...)
}

func (f *fakeKV) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	return f._do(ctx, op)
}

func (f *fakeKV) Txn(ctx context.Context) clientv3.Txn {
	return f._txn(ctx)
}
