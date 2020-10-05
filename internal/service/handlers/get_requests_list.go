package handlers

import (
	"net/http"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/service/handlers/models"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/service/requests"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetRequestsList(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetRequestsListRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	dataOwners := make([]string, 0)
	if request.FilterOwner != nil {
		dataOwners = append(dataOwners, *request.FilterOwner)
	}

	if !IsAllowed(r, w, dataOwners...) {
		return
	}

	requestsQ := RequestsQ(r).Page(request.OffsetPageParams)

	if request.FilterOwner != nil {
		requestsQ.FilterByOwner(*request.FilterOwner)
	}

	if len(request.FilterStatus) > 0 {
		statuses := make([]data.RequestStatus, len(request.FilterStatus))
		for i, status := range request.FilterStatus {
			statuses[i] = data.RequestStatus(status)
		}
		requestsQ.FilterByStatus(statuses...)
	}

	if request.FilterSourceBalance != nil {
		requestsQ.FilterBySourceBalance(*request.FilterSourceBalance)
	}

	if request.FilterAsset != nil {
		requestsQ.FilterByAsset(*request.FilterAsset)
	}

	requestsQ.FilterByCreatedAt(request.FilterFromCreatedAt, request.FilterToCreatedAt)

	requestsList, err := requestsQ.Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get requests from db")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	response := models.NewRequestListResponse(requestsList)
	response.Links = GetOffsetLinks(r, request.OffsetPageParams)
	ape.Render(w, response)
}
