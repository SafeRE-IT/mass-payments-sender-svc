package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetRequestsListRequest struct {
	pgdb.OffsetPageParams
	FilterOwner  *string  `filter:"owner"`
	FilterStatus []string `filter:"status"`
}

func NewGetRequestsListRequest(r *http.Request) (GetRequestsListRequest, error) {
	request := GetRequestsListRequest{}

	err := urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
