package httptool

import (
	"net/url"
	"path/filepath"
)

func UrlToFilepath(baseDir string, url *url.URL) string {
	// FIXME: query support
	return filepath.Join(baseDir, url.Scheme, url.Host, url.Path)
}
