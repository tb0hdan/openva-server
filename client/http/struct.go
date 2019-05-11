package httpclient

import "github.com/tb0hdan/openva-server/api"

type ForwarderLibrary struct {
	Items         []*api.LibraryItem `json:"items,omitempty"`
	StatusMessage string             `json:"status,omitempty"`
}
