package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type GetRequestRequest struct {
	ID int64 `url:"-"`
}

func NewGetRequestRequest(r *http.Request) (GetRequestRequest, error) {
	request := GetRequestRequest{}

	err := urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.ID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
