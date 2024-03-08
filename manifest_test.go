// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-10

package favicon

import (
	urls "net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParserAbsURL tests resolution of URLs.
func TestParserAbsURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		in, x string
		base  *urls.URL
	}{
		{"empty", "", "", nil},
		{"onlyBaseURL", "", "", mustURL("https://github.com")},
		{"noBaseURL", "/root", "/root", nil},
		{"baseURL", "/root", "https://github.com/root", mustURL("https://github.com")},
		{"absURL", "https://github.com/root", "https://github.com/root", mustURL("https://github.com")},
		// absolute URLs returned as-is
		{"absURLDifferentBase", "https://github.com/root", "https://github.com/root", mustURL("https://google.com")},
		{"absURLNoBase", "https://github.com/root", "https://github.com/root", nil},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			p := parser{baseURL: td.base}
			v := p.absURL(td.in)
			assert.Equal(t, td.x, v, "unexpected URL")
		})
	}
}

func mustURL(s string) *urls.URL {
	u, err := urls.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
