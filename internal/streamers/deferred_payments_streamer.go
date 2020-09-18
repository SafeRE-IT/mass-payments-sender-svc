package streamers

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"gitlab.com/distributed_lab/logan/v3/errors"

	"gitlab.com/distributed_lab/running"

	"gitlab.com/distributed_lab/json-api-connector/client"

	connector "gitlab.com/distributed_lab/json-api-connector"

	"gitlab.com/distributed_lab/logan/v3"

	regources "gitlab.com/tokend/regources/generated"
)

type DeferredPaymentsProcessor interface {
	ProcessPage(ctx context.Context, data regources.DeferredPaymentListResponse) error
}

type DeferredPaymentsStreamer struct {
	log       *logan.Entry
	streamer  *connector.Streamer
	processor DeferredPaymentsProcessor
}

func NewDeferredPaymentsStreamer(log *logan.Entry, client client.Client, params DeferredPaymentsQueryParams, processor DeferredPaymentsProcessor) *DeferredPaymentsStreamer {
	return &DeferredPaymentsStreamer{
		log:       log,
		streamer:  connector.NewStreamer(client, "v3/deferred_payments", params),
		processor: processor,
	}
}

func (s *DeferredPaymentsStreamer) Run(ctx context.Context) {
	running.WithBackOff(ctx, s.log, "deferred-payments-streamer", func(ctx context.Context) error {
		s.log.Debug("Try get next page")
		var page regources.DeferredPaymentListResponse
		err := s.streamer.Next(&page)
		if err != nil {
			return errors.Wrap(err, "failed to get next page of requests")
		}
		err = s.processor.ProcessPage(ctx, page)
		s.log.Debug("Finished processing page")
		if err != nil {
			return errors.Wrap(err, "failed to process request page")
		}

		return nil
	}, 30*time.Second, 30*time.Second, 30*time.Second)
}

type DeferredPaymentsQueryParams struct {
	Destination *string
	Cursor      *int64
}

func (p DeferredPaymentsQueryParams) Encode() string {
	values := url.Values{}

	if p.Destination != nil {
		values.Add("filter[destination]", *p.Destination)
	}

	if p.Cursor != nil {
		values.Add("page[cursor]", fmt.Sprintf("%d", *p.Cursor))
	}

	return values.Encode()
}
