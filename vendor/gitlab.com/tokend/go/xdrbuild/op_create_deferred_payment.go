package xdrbuild

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type CreateDeferredPayment struct {
	RequestID          uint64
	SourceBalance      string
	DestinationAccount string
	Amount             uint64
	Details            json.Marshaler
	AllTasks           *uint32
}

func (op *CreateDeferredPayment) XDR() (*xdr.Operation, error) {
	details, err := op.Details.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	var sb xdr.BalanceId
	err = sb.SetString(op.SourceBalance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set source balance id")
	}

	var da xdr.AccountId
	da = xdr.AccountId{}
	if err := da.SetAddress(op.DestinationAccount); err != nil {
		return nil, errors.Wrap(err, "failed to set op destination account id")
	}

	creationRequestOp := xdr.CreateDeferredPaymentCreationRequestOp{
		RequestId: xdr.Uint64(op.RequestID),
		Request: xdr.CreateDeferredPaymentRequest{
			SourceBalance:  sb,
			Destination:    da,
			Amount:         xdr.Uint64(op.Amount),
			CreatorDetails: xdr.Longstring(details),
			Ext:            xdr.EmptyExt{},
		},
		Ext: xdr.CreateDeferredPaymentCreationRequestOpExt{},
	}

	if op.AllTasks != nil {
		v := xdr.Uint32(*op.AllTasks)
		creationRequestOp.AllTasks = &v
	}

	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type:                                   xdr.OperationTypeCreateDeferredPaymentCreationRequest,
			CreateDeferredPaymentCreationRequestOp: &creationRequestOp,
		},
	}, nil
}
