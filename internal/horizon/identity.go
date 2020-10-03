package horizon

import (
	"encoding/json"
	"fmt"
	"net/url"

	"gitlab.com/tokend/mass-payments-sender-svc/resources"

	"github.com/pkg/errors"
)

func (c *Connector) GetIdentities(value, valueType string) (*IdentityData, error) {
	path, err := url.Parse(
		fmt.Sprintf("/integrations/identity-storage/identities?filter[value]=%s&filter[type]=%s?filter[status]=%s",
			value, valueType, "active"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse path")
	}
	var result IdentitiesResponse
	err = c.connector.Get(path, &result)

	if len(result.Data) == 0 {
		return nil, nil
	}

	return &result.Data[0], nil
}

type IdentitiesResponse struct {
	Data []IdentityData `json:"data"`
}

type IdentityData struct {
	Type          string                `json:"type"`
	ID            string                `json:"id"`
	Attributes    IdentityAttributes    `json:"attributes"`
	Relationships IdentityRelationships `json:"relationships"`
}

type IdentityRelationships struct {
	Account resources.Relation  `json:"account"`
	Type    resources.Relation  `json:"type"`
	Data    *resources.Relation `json:"data"`
}

type IdentityAttributes struct {
	Hash    string          `json:"hash"`
	Salt    string          `json:"salt"`
	Value   *string         `json:"value"`
	Details json.RawMessage `json:"details"`
	Status  string          `json:"status"`
}
