package service

import (
	"context"
	"net"
	"net/http"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/data/pg"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/streamers"

	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/config"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
	cfg      config.Config
}

func (s *service) run() error {
	s.runDeferredPaymentsStreamer()

	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func (s *service) runDeferredPaymentsStreamer() {
	processor := streamers.NewDeferredPaymentsProcessor(s.log, pg.NewRequestsQ(s.cfg.DB()))
	maxId, err := pg.NewRequestsQ(s.cfg.DB()).GetMaxId()
	if err != nil {
		panic(errors.Wrap(err, "failed to get max request id"))
	}
	dest := s.cfg.Keys().Source.Address()
	streamer := streamers.NewDeferredPaymentsStreamer(s.log, s.cfg.Client(),
		streamers.DeferredPaymentsQueryParams{
			Cursor:      maxId,
			Destination: &dest,
		}, processor)
	go streamer.Run(context.Background())
}

func newService(cfg config.Config) *service {
	return &service{
		log:      cfg.Log(),
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
		cfg:      cfg,
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}
