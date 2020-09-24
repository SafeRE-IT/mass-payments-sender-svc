package xdrbuild

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
)

type CloseDeferredPayment struct {
	RequestID          uint64
	DeferredPaymentID  uint64
	DestinationBalance string
	Amount             uint64
	FeeData            Fee
	Details            json.Marshaler
}

func (op *CloseDeferredPayment) XDR() (*xdr.Operation, error) {
	details, err := op.Details.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal details")
	}

	var db xdr.BalanceId
	db = xdr.BalanceId{}
	if err := db.SetString(op.DestinationBalance); err != nil {
		return nil, errors.Wrap(err, "failed to set op destination balance id")
	}

	return &xdr.Operation{
		Body: xdr.OperationBody{
			Type: xdr.OperationTypeCreateCloseDeferredPaymentRequest,
			CreateCloseDeferredPaymentRequestOp: &xdr.CreateCloseDeferredPaymentRequestOp{
				RequestId: xdr.Uint64(op.RequestID),
				Request: xdr.CloseDeferredPaymentRequest{
					DeferredPaymentId:  xdr.Uint64(op.DeferredPaymentID),
					DestinationBalance: xdr.BalanceId{},
					CreatorDetails:     xdr.Longstring(details),
					Amount:             xdr.Uint64(op.Amount),
					FeeData: xdr.PaymentFeeData{
						SourceFee: xdr.Fee{
							Fixed:   xdr.Uint64(op.FeeData.SourceFixed),
							Percent: xdr.Uint64(op.FeeData.SourcePercent),
						},
						DestinationFee: xdr.Fee{
							Fixed:   xdr.Uint64(op.FeeData.DestinationFixed),
							Percent: xdr.Uint64(op.FeeData.DestinationPercent),
						},
						SourcePaysForDest: op.FeeData.SourcePaysForDest,
					},
					Ext: xdr.EmptyExt{},
				},
				AllTasks: nil,
				Ext:      xdr.CreateCloseDeferredPaymentRequestOpExt{},
			},
		},
	}, nil
}
