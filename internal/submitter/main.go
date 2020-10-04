package submitter

import (
	"context"
	"time"

	"gitlab.com/tokend/keypair"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/tokend/go/xdrbuild"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/horizon"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"

	"gitlab.com/tokend/connectors/submit"

	"github.com/pkg/errors"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

func NewSubmitter(log *logan.Entry, paymentsQ data.PaymentsQ, horizonClient *horizon.Connector,
	signer keypair.Full, source keypair.Address) *Submitter {
	return &Submitter{
		log:           log,
		paymentsQ:     paymentsQ,
		horizonClient: horizonClient,
		signer:        signer,
		source:        source,
	}
}

type Submitter struct {
	log           *logan.Entry
	paymentsQ     data.PaymentsQ
	signer        keypair.Full
	source        keypair.Address
	horizonClient *horizon.Connector
}

func (s *Submitter) Run(ctx context.Context, submitPeriod int64, batchSize uint64) {
	period := time.Duration(submitPeriod)
	running.WithBackOff(ctx, s.log, "data-streamer", func(ctx context.Context) error {
		return s.processBatch(ctx, batchSize)
	}, period*time.Second, period*time.Second, period*time.Second)
}

func (s *Submitter) processBatch(ctx context.Context, size uint64) error {
	payments, err := s.paymentsQ.New().
		FilterByStatus(data.PaymentStatusProcessing).
		Limit(size).
		Select()
	if err != nil {
		return errors.Wrap(err, "failed to get txs from db")
	}

	for _, payment := range payments {
		payment := payment
		go func() {
			err := s.processTx(ctx, payment)
			if err != nil {
				s.log.WithError(err).
					WithField("payment_id", payment.ID).
					Warn("failed to submit transaction, try again later")
			}
		}()
	}

	return nil
}

func (s *Submitter) processTx(ctx context.Context, payment data.Payment) error {
	log := s.log.WithField("payment_id", payment.ID)
	log.Debug("sending transaction")
	defer log.Debug("payment sent")

	var err error
	if payment.TxBody == nil {
		payment.TxBody, err = s.buildCloseDeferredPaymentTx(payment)
		if err != nil {
			return errors.Wrap(err, "failed to build tx")
		}
		if payment.TxBody == nil {
			s.log.Infof("skiping payment %d", payment.ID)
			return nil
		}
		_, err = s.paymentsQ.SetTxBody(*payment.TxBody).FilterByID(payment.ID).Update()
	}
	_, err = s.horizonClient.Submit(ctx, *payment.TxBody, false)
	if err != nil {
		if err, ok := err.(*submit.TxFailure); ok {
			s.log.WithError(err).
				WithField("tx_hash", payment.ID).
				Warn("tx failed to submit marking it failed")
			_, err := s.paymentsQ.New().
				FilterByID(payment.ID).
				SetStatus(data.PaymentStatusFailed).
				SetFailureReason(err.Error()).
				Update()
			if err != nil {
				return errors.Wrap(err, "failed to mark transaction failed")
			}
		}
		return errors.Wrap(err, "failed to submit transaction")
	}

	return s.paymentsQ.New().Transaction(func(q data.PaymentsQ) error {
		_, err = q.FilterByID(payment.ID).
			SetStatus(data.PaymentStatusSuccess).
			Update()
		if err != nil {
			return errors.Wrap(err, "failed to mark transaction successes")
		}
		exists, err := q.Exists(payment.RequestID, data.PaymentStatusProcessing)
		if err != nil {
			return errors.Wrap(err, "failed to check processing transactions existence")
		}
		if exists {
			err = q.UpdateRequestStatus(payment.RequestID, data.RequestStatusFinished)
			if err != nil {
				return errors.Wrap(err, "failed to change request status")
			}
		}

		return nil
	})
}

func (s *Submitter) buildCloseDeferredPaymentTx(payment data.Payment) (*string, error) {
	if payment.DestinationType != data.DestinationTypeAccountID {
		identityData, err := s.horizonClient.GetIdentity(payment.Destination, payment.DestinationType)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get identities data")
		}
		if identityData == nil {
			s.log.Infof("identity not found for payment %d", payment.ID)
			return nil, nil
		}

		payment.Destination = identityData.Relationships.Account.Data.ID
		payment.DestinationType = data.DestinationTypeAccountID
	}

	horizonInfo, err := s.horizonClient.Info()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get horizon info")
	}

	var da xdr.AccountId
	err = da.SetAddress(payment.Destination)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get destination account id")
	}
	op := &xdrbuild.CloseDeferredPayment{
		Destination: xdr.CloseDeferredPaymentRequestDestination{
			Type:      xdr.CloseDeferredPaymentDestinationTypeAccount,
			AccountId: &da,
		},
		Amount:            uint64(payment.Amount),
		Details:           EmptyDetails{},
		DeferredPaymentID: uint64(payment.RequestID),
	}
	builder := xdrbuild.NewBuilder(horizonInfo.Attributes.NetworkPassphrase, horizonInfo.Attributes.TxExpirationPeriod)
	tx, err := builder.Transaction(s.source).Op(op).Sign(s.signer).Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal transaction")
	}

	return &tx, nil
}

type EmptyDetails struct{}

func (d EmptyDetails) MarshalJSON() ([]byte, error) {
	return []byte("{}"), nil
}