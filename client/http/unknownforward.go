package httpclient

import (
	"net/url"

	"github.com/tb0hdan/openva-server/api"
)

func UnknownForward(query, _ string) (items []*api.LibraryItem, err error) {
	fullURL := "http://localhost:49999/unknown/" + url.PathEscape(query)
	return Send(fullURL)
}
