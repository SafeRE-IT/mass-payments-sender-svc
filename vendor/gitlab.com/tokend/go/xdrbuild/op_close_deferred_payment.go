package xdrbuild

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type CloseDeferredPayment struct {
	RequestID         uint64
	DeferredPaymentID uint64
	Destination       xdr.CloseDeferredPaymentRequestDestination
	Amount            uint64
	Details           json.Marshaler
	AllTasks          *uint32
}

func (op *CloseDeferredPayment) XDR() (*xdr.Operation, error) {
	details, err := op.Details.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	requestOp := xdr.CreateCloseDeferredPaymentRequestOp{
		RequestId: xdr.Uint64(op.RequestID),
		Request: xdr.CloseDeferredPaymentRequest{
			DeferredPaymentId: xdr.Uint64(op.DeferredPaymentID),
			Destination:       op.Destination,
			CreatorDetails:    xdr.Longstring(details),
			Amount:            xdr.Uint64(op.Amount),
			Ext:               xdr.EmptyExt{},
		},
		Ext: xdr.CreateCloseDeferredPaymentRequestOpExt{},
	}

	if op.AllTasks != nil {
		v := xdr.Uint32(*op.AllTasks)
		requestOp.AllTasks = &v
	}

	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type:                                xdr.OperationTypeCreateCloseDeferredPaymentRequest,
			CreateCloseDeferredPaymentRequestOp: &requestOp,
		},
	}, nil
}
