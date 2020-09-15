package handlers

import (
	"net/http"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/service/handlers/models"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetRequest(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRequestRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	model, err := RequestsQ(r).FilterByID(request.ID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get request from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if model == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !IsAllowed(r, w, model.Owner) {
		return
	}

	ape.Render(w, models.NewRequestResponse(*model))
}
