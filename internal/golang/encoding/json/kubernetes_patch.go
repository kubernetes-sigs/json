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
	"strings"
)

// Type-alias error and data types returned from decoding

type UnmarshalTypeError = gojson.UnmarshalTypeError
type UnmarshalFieldError = gojson.UnmarshalFieldError
type InvalidUnmarshalError = gojson.InvalidUnmarshalError
type Number = gojson.Number
type RawMessage = gojson.RawMessage
type Token = gojson.Token
type Delim = gojson.Delim

type UnmarshalOpt func(*decodeState)

func UseNumber(d *decodeState) {
	d.useNumber = true
}
func DisallowUnknownFields(d *decodeState) {
	d.disallowUnknownFields = true
}

// saveStrictError saves a strict decoding error,
// for reporting at the end of the unmarshal if no other errors occurred.
func (d *decodeState) saveStrictError(err error) {
	// prevent excessive numbers of accumulated errors
	if len(d.savedStrictErrors) >= 100 {
		return
	}
	// dedupe accumulated strict errors
	if d.seenStrictErrors == nil {
		d.seenStrictErrors = map[string]struct{}{}
	}
	msg := err.Error()
	if _, seen := d.seenStrictErrors[msg]; seen {
		return
	}

	// accumulate the error
	d.seenStrictErrors[msg] = struct{}{}
	d.savedStrictErrors = append(d.savedStrictErrors, err)
}

// UnmarshalStrictError holds errors resulting from use of strict disallow___ decoder directives.
// If this is returned from Unmarshal(), it means the decoding was successful in all other respects.
type UnmarshalStrictError struct {
	Errors []error
}

func (e *UnmarshalStrictError) Error() string {
	var b strings.Builder
	b.WriteString("json: ")
	for i, err := range e.Errors {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}