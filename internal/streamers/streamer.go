package streamers

import (
	"context"
	"fmt"
	"net/url"
	"time"

	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/json-api-connector/client"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/running"

	connector "gitlab.com/distributed_lab/json-api-connector"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/connectors/request"
)

type Streamer interface {
	Run(ctx context.Context)
}

func NewCreateDeferredPayment(log *logan.Entry, processor request.PageProcessor, client client.Client,
	params RREncodeParams) Streamer {
	return NewBaseStreamer(log, processor, client, params, "v3/create_deferred_payments_requests")
}

func NewCreateCloseDeferredPayment(log *logan.Entry, processor request.PageProcessor, client client.Client,
	params RREncodeParams) Streamer {
	return NewBaseStreamer(log, processor, client, params, "v3/create_close_deferred_payments_requests")
}

type BaseStreamer struct {
	log       *logan.Entry
	processor request.PageProcessor
	connector *connector.Connector
	endpoint  *url.URL
}

func NewBaseStreamer(log *logan.Entry, processor request.PageProcessor, client client.Client,
	params RREncodeParams, rawEndpoint string) Streamer {
	endpoint, err := url.Parse(rawEndpoint)
	if err != nil {
		panic(errors.Wrap(err, fmt.Sprintf("failed to parse streamer url: %s", rawEndpoint)))
	}
	endpoint.RawQuery = params.Encode()

	return &BaseStreamer{
		log:       log,
		processor: processor,
		connector: connector.NewConnector(client),
		endpoint:  endpoint,
	}
}

func (s *BaseStreamer) Run(ctx context.Context) {
	running.WithBackOff(ctx, s.log, "request-streamer", func(ctx context.Context) error {
		var page regources.ReviewableRequestListResponse
		err := s.connector.Get(s.endpoint, &page)
		if err != nil {
			return errors.Wrap(err, "failed to get page of requests")
		}
		err = s.processor.ProcessPage(ctx, page)
		if err != nil {
			return errors.Wrap(err, "failed to process request page")
		}

		return nil
	}, 30*time.Second, 30*time.Second, 30*time.Second)
}

type RREncodeParams struct {
	State        int64
	PendingTasks uint32
	Destination  *string
}

func (p RREncodeParams) Encode() string {
	values := url.Values{}

	values.Add("filter[pending_tasks]", fmt.Sprintf("%d", p.PendingTasks))
	values.Add("filter[state]", fmt.Sprintf("%d", p.State))
	if p.Destination != nil {
		values.Add("filter[request_details.destination]", *p.Destination)
	}
	values.Add("include", "request_details")

	return values.Encode()
}
