package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"github.com/SafeRE-IT/mass-payments-sender-svc/resources"
)

func GetInfo(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, resources.InfoResponse{
		Data: resources.Info{
			Key: resources.NewKeyInt64(1, resources.MASS_PAYMENTS_INFO),
			Attributes: resources.InfoAttributes{
				AccountId: Keys(r).Address(),
			},
		},
	})
}
