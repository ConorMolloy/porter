package feed

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type MixinFeed struct {
	// Index of mixin files
	Index map[string]map[string]*MixinFileset

	// Mixins present in the feed
	Mixins []string

	// Updated timestamp according to the atom xml feed
	Updated *time.Time
}

func (feed *MixinFeed) Search(mixin string, version string) *MixinFileset {
	versions, ok := feed.Index[mixin]
	if !ok {
		return nil
	}

	return versions[version]
}

type MixinFileset struct {
	Mixin   string
	Version string
	Files   []MixinFile
}

func (f *MixinFileset) FindDownloadURL(os string, arch string) *url.URL {
	match := fmt.Sprintf("%s-%s-%s", f.Mixin, os, arch)
	for _, file := range f.Files {
		if strings.Contains(file.URL.Path, match) {
			return file.URL
		}
	}
	return nil
}

func (f *MixinFileset) Updated() string {
	return toAtomTimestamp(f.GetLastUpdated())
}

func (f *MixinFileset) GetLastUpdated() time.Time {
	var max time.Time
	for _, f := range f.Files {
		if f.Updated.After(max) {
			max = f.Updated
		}
	}
	return max
}

type MixinFile struct {
	File    string
	URL     *url.URL
	Updated time.Time
}

// MixinEntries is used to sort the entries in a mixin feed by when they were last updated
type MixinEntries []*MixinFileset

func (e MixinEntries) Len() int {
	return len(e)
}

func (e MixinEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e MixinEntries) Less(i, j int) bool {
	return e[i].GetLastUpdated().Before(e[j].GetLastUpdated())
}