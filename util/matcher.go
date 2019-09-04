package util

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/dnaeon/go-vcr/cassette"
)

// SetMatcher sets matcher for VCR used in tests
func SetMatcher(r *http.Request, i cassette.Request) bool {
	var b bytes.Buffer
	if _, err := b.ReadFrom(r.Body); err != nil {
		return false
	}
	r.Body = ioutil.NopCloser(&b)
	return cassette.DefaultMatcher(r, i) && (b.String() == "" || b.String() == i.Body)
}
