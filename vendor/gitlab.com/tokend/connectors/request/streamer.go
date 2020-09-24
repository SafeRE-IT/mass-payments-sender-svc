package request

import (
	"context"
	connector "gitlab.com/distributed_lab/json-api-connector"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/running"
	regources "gitlab.com/tokend/regources/generated"
	"time"
)

type PageProcessor interface {
	ProcessPage(context.Context, regources.ReviewableRequestListResponse) error
}

type Streamer struct {
	log       *logan.Entry
	processor PageProcessor
	streamer  *connector.Streamer
}

type StreamerOpts struct {
	Log       *logan.Entry
	Processor PageProcessor
	Streamer  *connector.Streamer
}

func NewStreamer(opts StreamerOpts) *Streamer {
	return &Streamer{
		log:       opts.Log,
		streamer:  opts.Streamer,
		processor: opts.Processor,
	}
}

func (s *Streamer) Run(ctx context.Context) {
	running.WithBackOff(ctx, s.log, "request-streamer", func(ctx context.Context) error {
		var page regources.ReviewableRequestListResponse
		err := s.streamer.Next(&page)
		if err != nil {
			return errors.Wrap(err, "failed to get next page of requests")
		}
		err = s.processor.ProcessPage(ctx, page)
		if err != nil {
			return errors.Wrap(err, "failed to process request page")
		}

		return nil
	}, 30*time.Second, 30*time.Second, 10*time.Minute)
}
