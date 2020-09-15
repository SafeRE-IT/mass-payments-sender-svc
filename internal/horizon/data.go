package horizon

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	regources "gitlab.com/tokend/regources/generated"
)

func (c *Connector) DataList(dataType int) (regources.DataListResponse, error) {
	endpoint, err := url.Parse(fmt.Sprintf("/v3/data?filter[type]=%d", dataType))
	if err != nil {
		return regources.DataListResponse{}, errors.Wrap(err, "failed to parse url")
	}

	var response regources.DataListResponse
	err = c.connector.Get(endpoint, &response)
	if err != nil {
		return response, errors.Wrap(err, "failed to get key value")
	}

	return response, nil
}
