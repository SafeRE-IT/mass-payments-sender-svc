package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetPaymentsListRequest struct {
	pgdb.OffsetPageParams
	FilterRequestId *int64   `filter:"request_id"`
	FilterStatus    []string `filter:"status"`
}

func NewGetPaymentsListRequest(r *http.Request) (GetPaymentsListRequest, error) {
	request := GetPaymentsListRequest{}

	err := urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
