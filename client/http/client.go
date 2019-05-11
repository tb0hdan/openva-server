package httpclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/tb0hdan/openva-server/api"
)

func Send(fullURL string) (items []*api.LibraryItem, err error) {
	resp, err := http.DefaultClient.Post(fullURL, "", nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	fwd := &ForwarderLibrary{}
	err = json.Unmarshal(data, fwd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if fwd.Items != nil {
		return fwd.Items, nil
	}
	return
}
