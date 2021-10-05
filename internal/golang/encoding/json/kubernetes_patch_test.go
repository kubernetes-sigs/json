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
	"strings"
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

func TestStrictErrors(t *testing.T) {
	type Typed struct {
		A int `json:"a"`
	}

	testcases := []struct {
		name            string
		in              string
		expectStrictErr bool
		expectErr       string
	}{
		{
			name:            "malformed 1",
			in:              `{`,
			expectStrictErr: false,
		},
		{
			name:            "malformed 2",
			in:              `{}}`,
			expectStrictErr: false,
		},
		{
			name:            "malformed 3",
			in:              `{,}`,
			expectStrictErr: false,
		},
		{
			name:            "type error",
			in:              `{"a":true}`,
			expectStrictErr: false,
		},
		{
			name:            "unknown",
			in:              `{"a":1,"unknown":true,"unknown":false}`,
			expectStrictErr: true,
			expectErr:       `json: unknown field "unknown"`,
		},
		{
			name:            "unknowns",
			in:              `{"a":1,"unknown":true,"unknown2":true,"unknown":true,"unknown2":true}`,
			expectStrictErr: true,
			expectErr:       `json: unknown field "unknown", unknown field "unknown2"`,
		},
		{
			name:            "unknowns and type error",
			in:              `{"unknown":true,"a":true}`,
			expectStrictErr: false,
		},
		{
			name:            "unknowns and malformed error",
			in:              `{"unknown":true}}`,
			expectStrictErr: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := Unmarshal([]byte(tc.in), &Typed{}, DisallowUnknownFields)
			if err == nil {
				t.Fatal("expected error, got none")
			}
			_, isStrictErr := err.(*UnmarshalStrictError)
			if tc.expectStrictErr != isStrictErr {
				t.Fatalf("expected strictErr=%v, got %v: %v", tc.expectStrictErr, isStrictErr, err)
			}
			if !strings.Contains(err.Error(), tc.expectErr) {
				t.Fatalf("expected error containing %q, got %q", tc.expectErr, err)
			}
			t.Log(err)
		})
	}
}

func TestCaseSensitive(t *testing.T) {
	type Embedded1 struct {
		C int `json:"c"`
		D int
	}
	type Embedded2 struct {
		E int `json:"e"`
		F int
	}

	type Obj struct {
		A         int `json:"a"`
		B         int
		Embedded1 `json:",inline"`
		Embedded2
	}

	testcases := []struct {
		name   string
		in     string
		to     interface{}
		expect interface{}
	}{
		{
			name:   "tagged",
			in:     `{"A":"1","A":2,"a":3,"A":4,"A":"5"}`,
			to:     &Obj{},
			expect: &Obj{A: 3},
		},
		{
			name:   "untagged",
			in:     `{"b":"1","b":2,"B":3,"b":4,"b":"5"}`,
			to:     &Obj{},
			expect: &Obj{B: 3},
		},
		{
			name:   "inline embedded tagged subfield",
			in:     `{"C":"1","C":2,"c":3,"C":4,"C":"5"}`,
			to:     &Obj{},
			expect: &Obj{Embedded1: Embedded1{C: 3}},
		},
		{
			name:   "inline embedded untagged subfield",
			in:     `{"d":"1","d":2,"D":3,"d":4,"d":"5"}`,
			to:     &Obj{},
			expect: &Obj{Embedded1: Embedded1{D: 3}},
		},
		{
			name:   "inline embedded field name",
			in:     `{"Embedded1":{"c":3}}`,
			to:     &Obj{},
			expect: &Obj{}, // inlined embedded is not addressable by field name
		},
		{
			name:   "inline embedded empty name",
			in:     `{"":{"c":3}}`,
			to:     &Obj{},
			expect: &Obj{}, // inlined embedded is not addressable by empty json field name
		},
		{
			name:   "untagged embedded tagged subfield",
			in:     `{"E":"1","E":2,"e":3,"E":4,"E":"5"}`,
			to:     &Obj{},
			expect: &Obj{Embedded2: Embedded2{E: 3}},
		},
		{
			name:   "untagged embedded untagged subfield",
			in:     `{"f":"1","f":2,"F":3,"f":4,"f":"5"}`,
			to:     &Obj{},
			expect: &Obj{Embedded2: Embedded2{F: 3}},
		},
		{
			name:   "untagged embedded field name",
			in:     `{"Embedded2":{"e":3}}`,
			to:     &Obj{},
			expect: &Obj{}, // untagged embedded is not addressable by field name
		},
		{
			name:   "untagged embedded empty name",
			in:     `{"":{"e":3}}`,
			to:     &Obj{},
			expect: &Obj{}, // untagged embedded is not addressable by empty json field name
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if err := Unmarshal([]byte(tc.in), &tc.to, CaseSensitive); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tc.expect, tc.to) {
				t.Fatalf("expected\n%#v\ngot\n%#v", tc.expect, tc.to)
			}
		})
	}
}
