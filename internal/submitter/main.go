package submitter

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/tokend/go/xdr"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/horizon"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"

	"gitlab.com/tokend/connectors/submit"

	"github.com/pkg/errors"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

func NewSubmitter(log *logan.Entry, paymentsQ data.PaymentsQ, horizonClient *horizon.Connector) *Submitter {
	return &Submitter{
		log:           log,
		paymentsQ:     paymentsQ,
		horizonClient: horizonClient,
	}
}

type Submitter struct {
	log           *logan.Entry
	paymentsQ     data.PaymentsQ
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

	err := s.sendPayment(ctx, payment)
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

func (s *Submitter) sendPayment(ctx context.Context, payment data.Payment) error {
	s.log.Info(fmt.Sprintf("sending payment %d", payment.ID))

	// Get tx from db
	// If has send

	// Get account_id
	// If hasn't skip

	// Build tx

	// Save to db

	// Send to core
	return nil
}

func (s *Submitter) buildCloseDeferredPaymentTx(payment data.Payment) (xdr.TransactionEnvelope, error) {
	return xdr.TransactionEnvelope{}, nil
}
