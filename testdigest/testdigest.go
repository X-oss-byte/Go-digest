// Copyright 2021 OCI Contributors
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

// testdigest is a separate package, because it has some testing utilities in it that may be useful
// to other internal Algorithm implementors.
//
// It is not a stable interface and not meant for consumption outside of digest developers.

package testdigest

import (
	"testing"

	pkgdigest "github.com/opencontainers/go-digest"
)

type TestCase struct {
	// Input the formal format of the hash, for example sha256:e58fcf7418d4390dec8e8fb69d88c06ec07039d651fedd3aa72af9972e7d046b
	Input string
	// If err is non-nil, then the parsing of Input is expected to return this error
	Err error
	// Algorithm should be an available or registered algorithm
	Algorithm pkgdigest.Algorithm
	// Encoded is the the encoded portion of the digest to expect, for example e58fcf7418d4390dec8e8fb69d88c06ec07039d651fedd3aa72af9972e7d046b
	Encoded string
}

func RunTestCase(t *testing.T, testcase TestCase) {
	digest, err := pkgdigest.Parse(testcase.Input)
	if err != testcase.Err {
		t.Fatalf("error differed from expected while parsing %q: %v != %v", testcase.Input, err, testcase.Err)
	}

	if testcase.Err != nil {
		return
	}

	if digest.Algorithm() != testcase.Algorithm {
		t.Fatalf("incorrect Algorithm for parsed digest: %q != %q", digest.Algorithm(), testcase.Algorithm)
	}

	if digest.Encoded() != testcase.Encoded {
		t.Fatalf("incorrect hex for parsed digest: %q != %q", digest.Encoded(), testcase.Encoded)
	}

	// Parse string return value and check equality
	newParsed, err := pkgdigest.Parse(digest.String())
	if err != nil {
		t.Fatalf("unexpected error parsing Input %q: %v", testcase.Input, err)
	}

	if newParsed != digest {
		t.Fatalf("expected equal: %q != %q", newParsed, digest)
	}

	newFromHex := pkgdigest.NewDigestFromEncoded(newParsed.Algorithm(), newParsed.Encoded())
	if newFromHex != digest {
		t.Fatalf("%v != %v", newFromHex, digest)
	}
}
