package streamers

import (
	"context"
	"fmt"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/data"
	regources "gitlab.com/tokend/regources/generated"
)

type deferredPaymentsProcessor struct {
	log       *logan.Entry
	requestsQ data.RequestsQ
}

func NewDeferredPaymentsProcessor(log *logan.Entry, requestsQ data.RequestsQ) DeferredPaymentsProcessor {
	return &deferredPaymentsProcessor{
		log:       log,
		requestsQ: requestsQ,
	}
}

func (p *deferredPaymentsProcessor) ProcessPage(ctx context.Context, page regources.DeferredPaymentListResponse) error {
	for _, payment := range page.Data {
		p.log.Debug(fmt.Sprintf("Get deferred payment %s", payment.ID))
	}

	return nil
}
