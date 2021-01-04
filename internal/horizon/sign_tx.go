package horizon

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

func (c *Connector) SignTx(txBase64 string) (string, error) {
	path, err := url.Parse("/integrations/mass-payments/sign")
	if err != nil {
		return "", errors.Wrap(err, "failed to parse path")
	}
	var result tx
	err = c.connector.PostJSON(path, tx{Tx: txBase64}, context.TODO(), &result)
	if err != nil {
		return "", errors.Wrap(err, "failed to send sign request")
	}

	return result.Tx, nil
}

type tx struct {
	Tx string `json:"tx"`
}
