package horizon

import (
	"fmt"
	"net/url"

	regources "gitlab.com/tokend/regources/generated"

	"github.com/pkg/errors"
)

func (c *Connector) GetBalance(id string) (*regources.BalanceResponse, error) {
	endpoint, err := url.Parse(fmt.Sprintf("/v3/balances/%s", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	var response regources.BalanceResponse
	err = c.connector.Get(endpoint, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get blob")
	}

	return &response, nil
}
