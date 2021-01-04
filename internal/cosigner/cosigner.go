package cosigner

import (
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"gitlab.com/tokend/connectors/signed"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/config"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/horizon"
)

type Cosigner interface {
	Cosign(txBase64 string) (string, error)
}

func NewCosigner(config config.DecentralizationConfig) Cosigner {
	connectors := make([]*horizon.Connector, len(config.Nodes))
	for i, rawEndpoint := range config.Nodes {
		endpoint, err := url.Parse(rawEndpoint)
		if err != nil {
			panic(err)
		}
		connectors[i] = horizon.NewConnector(signed.NewClient(http.DefaultClient, endpoint))
	}

	return &cosigner{
		disabled:   config.Disabled,
		connectors: connectors,
	}
}

type cosigner struct {
	disabled   bool
	connectors []*horizon.Connector
}

func (c *cosigner) Cosign(txBase64 string) (string, error) {
	if c.disabled {
		return txBase64, nil
	}

	var err error
	for _, connector := range c.connectors {
		txBase64, err = connector.SignTx(txBase64)
		if err != nil {
			return "", errors.Wrap(err, "failed to cosign transaction")
		}
	}

	return txBase64, nil
}
