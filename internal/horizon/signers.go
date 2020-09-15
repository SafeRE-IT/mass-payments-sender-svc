package horizon

import (
	"fmt"
	"net/url"

	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/tokend/go/resources"
)

func (c *Connector) Signers(address string) ([]resources.Signer, error) {
	path, err := url.Parse(fmt.Sprintf("/v3/accounts/%s/signers", address))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse url")
	}

	var res regources.SignerListResponse
	err = c.connector.Get(path, &res)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get signers from horizon")
	}

	signers := make([]resources.Signer, 0, len(res.Data))
	for _, raw := range res.Data {
		signers = append(signers, resources.Signer{
			AccountID: raw.ID,
			Weight:    cast.ToInt(raw.Attributes.Weight),
			Identity:  raw.Attributes.Identity,
		})
	}

	return signers, nil
}
