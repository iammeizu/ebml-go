// Copyright 2019 The ebml-go authors.
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

package ebml

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMarshal_Omitempty(t *testing.T) {
	type TestOmitempty struct {
		EBML struct {
			DocType        string `ebml:"EBMLDocType,omitempty"`
			DocTypeVersion uint64 `ebml:"EBMLDocTypeVersion,omitempty"`
		} `ebml:"EBML"`
	}
	type TestNoOmitempty struct {
		EBML struct {
			DocType        string `ebml:"EBMLDocType"`
			DocTypeVersion uint64 `ebml:"EBMLDocTypeVersion"`
		} `ebml:"EBML"`
	}

	testCases := map[string]struct {
		input    interface{}
		expected []byte
	}{
		"Omitempty": {
			&TestOmitempty{},
			[]byte{0x1a, 0x45, 0xDF, 0xA3, 0x80},
		},
		"NoOmitempty": {
			&TestNoOmitempty{},
			[]byte{0x1A, 0x45, 0xDF, 0xA3, 0x88, 0x42, 0x82, 0x81, 0x00, 0x42, 0x87, 0x81, 0x00},
		},
	}

	for n, c := range testCases {
		t.Run(n, func(t *testing.T) {
			var b bytes.Buffer
			if err := Marshal(c.input, &b); err != nil {
				t.Fatalf("error: %+v\n", err)
			}
			if bytes.Compare(c.expected, b.Bytes()) != 0 {
				t.Errorf("Marshaled binary doesn't match:\n expected: %v,\n      got: %v", c.expected, b.Bytes())
			}
		})
	}
}

func ExampleMarshal() {
	type EBMLHeader struct {
		DocType            string `ebml:"EBMLDocType"`
		DocTypeVersion     uint64 `ebml:"EBMLDocTypeVersion"`
		DocTypeReadVersion uint64 `ebml:"EBMLDocTypeReadVersion"`
	}
	type TestEBML struct {
		Header EBMLHeader `ebml:"EBML"`
	}
	s := TestEBML{
		Header: EBMLHeader{
			DocType:            "webm",
			DocTypeVersion:     2,
			DocTypeReadVersion: 2,
		},
	}

	var b bytes.Buffer
	if err := Marshal(&s, &b); err != nil {
		panic(err)
	}
	for _, b := range b.Bytes() {
		fmt.Printf("0x%02x, ", int(b))
	}
	// Output:
	// 0x1a, 0x45, 0xdf, 0xa3, 0x90, 0x42, 0x82, 0x85, 0x77, 0x65, 0x62, 0x6d, 0x00, 0x42, 0x87, 0x81, 0x02, 0x42, 0x85, 0x81, 0x02,
}

func TestMarshal_Tag(t *testing.T) {
	tagged := struct {
		DocCustomNamedType string `ebml:"EBMLDocType"`
	}{
		DocCustomNamedType: "hoge",
	}
	untagged := struct {
		EBMLDocType string
	}{
		EBMLDocType: "hoge",
	}

	var bTagged, bUntagged bytes.Buffer
	if err := Marshal(&tagged, &bTagged); err != nil {
		t.Fatalf("error: %+v\n", err)
	}
	if err := Marshal(&untagged, &bUntagged); err != nil {
		t.Fatalf("error: %+v\n", err)
	}

	if bytes.Compare(bTagged.Bytes(), bUntagged.Bytes()) != 0 {
		t.Errorf("Tagged struct and untagged struct must be marshal-ed to same binary, tagged: %v, untagged: %v", bTagged.Bytes(), bUntagged.Bytes())
	}
}

func BenchmarkMarshal(b *testing.B) {
	type EBMLHeader struct {
		DocType            string `ebml:"EBMLDocType"`
		DocTypeVersion     uint64 `ebml:"EBMLDocTypeVersion"`
		DocTypeReadVersion uint64 `ebml:"EBMLDocTypeReadVersion"`
	}
	type TestEBML struct {
		Header EBMLHeader `ebml:"EBML"`
	}
	s := TestEBML{
		Header: EBMLHeader{
			DocType:            "webm",
			DocTypeVersion:     2,
			DocTypeReadVersion: 2,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		if err := Marshal(&s, &buf); err != nil {
			b.Fatalf("error: %+v\n", err)
		}
	}
}
