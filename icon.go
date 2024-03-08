// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-11-09

package favicon

import (
	"fmt"
	"sort"
)

// Icon is a favicon parsed from an HTML file or JSON manifest.
//
// TODO: Use *Icon everywhere to be consistent with higher-level APIs that return nil for "not found".
type Icon struct {
	URL      string `json:"url"`       // Never empty
	MimeType string `json:"mimetype"`  // MIME type of icon; never empty
	FileExt  string `json:"extension"` // File extension; may be empty
}

// String implements Stringer.
func (i Icon) String() string {
	return fmt.Sprintf("Icon{URL: %q,MimeType: %q}", i.URL, i.MimeType)
}

// Copy returns a new Icon with the same values as this one.
func (i Icon) Copy() *Icon {
	return &Icon{
		URL:      i.URL,
		MimeType: i.MimeType,
		FileExt:  i.FileExt,
	}
}

type ByURL []*Icon

// Implement sort.Interface
func (v ByURL) Len() int      { return len(v) }
func (v ByURL) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

func (v ByURL) Less(i, j int) bool {
	a, b := v[i], v[j]
	if a == nil || b == nil {
		return true
	}
	return a.URL < b.URL
}

// Check missing values, remove duplicates, sort.
func (p *parser) postProcessIcons(icons []*Icon) []*Icon {
	tidied := map[string]*Icon{}
	for _, icon := range icons {
		icon.URL = p.absURL(icon.URL)

		if icon.MimeType == "" {
			icon.MimeType = mimeTypeURL(icon.URL)
		}

		if icon.URL == "" || icon.MimeType == "" {
			continue
		}

		if icon.FileExt == "" {
			icon.FileExt = fileExt(icon.URL)
		}
		tidied[icon.URL] = icon
	}

	icons = []*Icon{}
	for _, icon := range tidied {
		for _, fun := range p.find.filters {
			if icon = fun(icon); icon == nil {
				break
			}
		}
		if icon != nil {
			icons = append(icons, icon)
		}
	}

	sort.Sort(ByURL(icons))
	return icons
}
