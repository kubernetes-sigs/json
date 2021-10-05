/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package json

import (
	gojson "encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalWithOptions(t *testing.T) {
	type Typed struct {
		A int `json:"a"`
	}

	testcases := []struct {
		name      string
		in        string
		to        interface{}
		options   []UnmarshalOpt
		expect    interface{}
		expectErr bool
	}{
		{
			name:   "default untyped",
			in:     `{"a":1}`,
			to:     map[string]interface{}{},
			expect: map[string]interface{}{"a": float64(1)},
		},
		{
			name:   "default typed",
			in:     `{"a":1, "unknown":"foo"}`,
			to:     &Typed{},
			expect: &Typed{A: 1},
		},
		{
			name:    "usenumbers untyped",
			in:      `{"a":1}`,
			to:      map[string]interface{}{},
			options: []UnmarshalOpt{UseNumber},
			expect:  map[string]interface{}{"a": gojson.Number("1")},
		},
		{
			name:   "usenumbers typed",
			in:     `{"a":1}`,
			to:     &Typed{},
			expect: &Typed{A: 1},
		},
		{
			name:    "disallowunknown untyped",
			in:      `{"a":1,"unknown":"foo"}`,
			to:      map[string]interface{}{},
			options: []UnmarshalOpt{DisallowUnknownFields},
			expect:  map[string]interface{}{"a": float64(1), "unknown": "foo"},
		},
		{
			name:      "disallowunknown typed",
			in:        `{"a":1,"unknown":"foo"}`,
			to:        &Typed{},
			options:   []UnmarshalOpt{DisallowUnknownFields},
			expectErr: true,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := Unmarshal([]byte(tc.in), &tc.to, tc.options...)
			if tc.expectErr != (err != nil) {
				t.Fatalf("expected err=%v, got %v", tc.expectErr, err)
			}
			if tc.expectErr {
				return
			}
			if !reflect.DeepEqual(tc.expect, tc.to) {
				t.Fatalf("expected\n%#v\ngot\n%#v", tc.expect, tc.to)
			}
		})
	}
}
