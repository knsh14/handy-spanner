// Copyright 2019 Masahiro Sano
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"testing"

	cmp "github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func TestQueryBuilder_UnnestValue_Success(t *testing.T) {
	b := newTestQueryBuilder()

	table := []struct {
		v    Value
		ph   string
		args []interface{}
	}{
		{
			v: Value{
				Data: []bool{true, false},
			},
			ph:   "VALUES (?), (?)",
			args: []interface{}{true, false},
		},
		{
			v: Value{
				Data: []int64{(100), int64(101)},
			},
			ph:   "VALUES (?), (?)",
			args: []interface{}{int64(100), int64(101)},
		},
		{
			v: Value{
				Data: []float64{float64(1.1), float64(1.2)},
			},
			ph:   "VALUES (?), (?)",
			args: []interface{}{float64(1.1), float64(1.2)},
		},
		{
			v: Value{
				Data: []string{"aa", "bb", "cc"},
			},
			ph:   "VALUES (?), (?), (?)",
			args: []interface{}{"aa", "bb", "cc"},
		},
		{
			v: Value{
				Data: [][]byte{[]byte("aa"), []byte("bb"), []byte("cc")},
			},
			ph:   "VALUES (?), (?), (?)",
			args: []interface{}{[]byte("aa"), []byte("bb"), []byte("cc")},
		},
	}

	for _, tc := range table {
		s, d, err := b.unnestValue(tc.v, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if s.Raw != tc.ph {
			t.Errorf("expect placeholder %q, but got %q", tc.ph, s.Raw)
		}

		if diff := cmp.Diff(tc.args, d); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}

	}
}

func TestQueryBuilder_UnnestValue_Error(t *testing.T) {
	b := newTestQueryBuilder()

	table := []struct {
		v    Value
		code codes.Code
	}{
		{
			v: Value{
				Data: nil,
			},
			code: codes.Unknown,
		},
		{
			v: Value{
				Data: true,
			},
			code: codes.InvalidArgument,
		},
		{
			v: Value{
				Data: int64(100),
			},
			code: codes.InvalidArgument,
		},
		{
			v: Value{
				Data: float64(100),
			},
			code: codes.InvalidArgument,
		},
		{
			v: Value{
				Data: "x",
			},
			code: codes.InvalidArgument,
		},
		{
			v: Value{
				Data: []byte("x"),
			},
			code: codes.InvalidArgument,
		},
	}

	for _, tc := range table {
		_, _, err := b.unnestValue(tc.v, false)
		st := status.Convert(err)
		if st.Code() != tc.code {
			t.Errorf("expect code %v, but got %v", tc.code, st.Code())
		}
	}
}
