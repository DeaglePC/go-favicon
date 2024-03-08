// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-09

package favicon

import (
	"encoding/json"
	"io"
)

// Manifest is the relevant parts of a manifest.json file.
type Manifest struct {
	Icons []ManifestIcon `json:"icons"`
}

// ManifestIcon is an icon from a manifest.json file.
type ManifestIcon struct {
	URL      string `json:"src"`
	Type     string `json:"type"`
	RawSizes string `json:"sizes"`
}

func (p *parser) parseManifest(url string) []*Icon {
	rc, err := p.find.fetchURL(url)
	if err != nil {
		return nil
	}
	defer rc.Close()

	return p.parseManifestReader(rc)
}

func (p *parser) parseManifestReader(r io.Reader) []*Icon {
	var (
		icons []*Icon
		man   = Manifest{}
		err   error
	)

	dec := json.NewDecoder(r)
	if err = dec.Decode(&man); err != nil {
	}
	for _, mi := range man.Icons {
		// TODO: make URL relative to manifest, not page
		mi.URL = p.absURL(mi.URL)
		icon := &Icon{URL: mi.URL}
		icons = append(icons, icon)
	}

	return icons
}
