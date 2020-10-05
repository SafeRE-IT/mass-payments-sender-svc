package requests

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetRequestsListRequest struct {
	pgdb.OffsetPageParams
	FilterOwner         *string    `filter:"owner"`
	FilterStatus        []string   `filter:"status"`
	FilterAsset         *string    `filter:"asset"`
	FilterSourceBalance *string    `filter:"source_balance"`
	FilterFromCreatedAt *time.Time `filter:"from_created_at"`
	FilterToCreatedAt   *time.Time `filter:"to_created_at"`
}

func NewGetRequestsListRequest(r *http.Request) (GetRequestsListRequest, error) {
	request := GetRequestsListRequest{}

	err := urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
