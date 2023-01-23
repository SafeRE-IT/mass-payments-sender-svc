package handlers

import (
	"net/http"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"github.com/SafeRE-IT/mass-payments-sender-svc/internal/service/requests"
)

func SignTx(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewSignTxRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !IsAllowed(r, w) {
		return
	}

	var txEnvelope xdr.TransactionEnvelope
	if err := xdr.SafeUnmarshalBase64(request.TxBase64, &txEnvelope); err != nil {
		ape.RenderErr(w, problems.BadRequest(
			validation.Errors{"tx": errors.New("should be valid base64 encoded transaction envelope")})...)
		return
	}

	resultTx, err := XdrBuilder(r).Sign(&txEnvelope, Signer(r))
	if err != nil {
		Log(r).WithError(err).Error("failed to cosign tx")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	base64, err := xdr.MarshalBase64(resultTx)
	if err != nil {
		Log(r).WithError(err).Error("failed to encode tx")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, requests.Tx{
		TxBase64: base64,
	})
}
