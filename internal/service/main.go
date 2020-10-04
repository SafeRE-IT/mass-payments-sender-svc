package service

import (
	"context"
	"net"
	"net/http"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/submitter"

	"gitlab.com/tokend/connectors/request"
	"gitlab.com/tokend/mass-payments-sender-svc/internal/horizon"

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
	s.runSubmitter()

	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func (s *service) runSubmitter() {
	horizonClient := horizon.NewConnector(s.cfg.Client())
	go submitter.
		NewSubmitter(s.log, pg.NewPaymentsQ(s.cfg.DB()), horizonClient, s.cfg.Keys().Signer, s.cfg.Keys().Source).
		Run(context.Background(), s.cfg.MassPaymentsSenderConfig().TxsPerPeriod, uint64(s.cfg.MassPaymentsSenderConfig().TxsPerPeriod))
}

func (s *service) runDeferredPaymentsStreamer() {
	processor := streamers.NewCreateDeferredPaymentRequestProcessor(s.log, pg.NewRequestsQ(s.cfg.DB()),
		request.NewReviewer(s.cfg), horizon.NewConnector(s.cfg.Client()), 1)
	dest := s.cfg.Keys().Source.Address()
	streamer := streamers.NewCreateDeferredPaymentStreamer(s.cfg.Log(), processor,
		s.cfg.Client(), streamers.RREncodeParams{
			State:        1,
			PendingTasks: 1,
			Destination:  &dest,
		})
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
