package dbbyte

import (
	"net/url"
	"strings"
)

// encodeKey encodes an object key for use in a URL path. Object keys can
// contain slashes (e.g. "images/logo.png"). The API's catch-all routes
// expect each path segment encoded individually, then rejoined with "/" -
// encoding the whole key would also encode the slashes themselves.
func encodeKey(key string) string {
	segments := strings.Split(key, "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}
