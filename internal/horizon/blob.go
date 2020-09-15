package horizon

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

func (c *Connector) GetBlob(id string) (*BlobResponse, error) {
	endpoint, err := url.Parse(fmt.Sprintf("/blobs/%s", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	var response BlobResponse
	err = c.connector.Get(endpoint, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get blob")
	}

	return &response, nil
}

type Blob struct {
	Attributes BlobAttributes `json:"attributes"`
}

type BlobAttributes struct {
	Value string `json:"value"`
}

type BlobResponse struct {
	Data     Blob     `json:"data"`
}
