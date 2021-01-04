package requests

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Tx struct {
	TxBase64 string `json:"tx"`
}

func NewSignTxRequest(r *http.Request) (Tx, error) {
	request := Tx{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, nil
}
