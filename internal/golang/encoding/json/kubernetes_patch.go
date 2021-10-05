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

import gojson "encoding/json"

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
