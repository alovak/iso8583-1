// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package fuzzreader

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/moov-io/iso8583/pkg/lib"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

// Return codes (from go-fuzz docs)
//
// The function must return 1 if the fuzzer should increase priority
// of the given input during subsequent fuzzing (for example, the input is
// lexically correct and was parsed successfully); -1 if the input must not be
// added to corpus even if gives new coverage; and 0 otherwise; other values are
// reserved for future use.
func Fuzz(data []byte) int {
	jsonData, err := ioutil.ReadFile(filepath.Join(basePath, "..", "testdata", "specification_ver_1987.json"))
	if err != nil {
		return -1
	}

	spec, err := lib.NewSpecificationWithJson(jsonData)
	if err != nil {
		return -1
	}

	message, err := lib.NewISO8583Message(spec)
	if err != nil {
		return -1
	}

	// Parse from raw data
	read, err := message.Load(data)
	if err != nil {
		return 0
	}

	// Check read size
	if read != len(data) {
		return 0
	}

	// Validate message
	err = message.Validate()
	if err != nil {
		return 0
	}

	return 1
}
