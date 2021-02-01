package service

import (
	"context"
	"net"
	"net/http"

	"gitlab.com/tokend/go/xdrbuild"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/cosigner"

	"gitlab.com/tokend/mass-payments-sender-svc/internal/submitter"

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
	if s.cfg.DecentralizationConfig().IsCoordinator || s.cfg.DecentralizationConfig().Disabled {
		s.runDeferredPaymentsStreamer()
		s.runSubmitter()
	}

	r := s.router()

	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func (s *service) runSubmitter() {
	horizonClient := horizon.NewConnector(s.cfg.Client())
	go submitter.
		NewSubmitter(s.log, pg.NewPaymentsQ(s.cfg.DB()), pg.NewRequestsQ(s.cfg.DB()), horizonClient, s.cfg.Keys().Signer,
			s.cfg.Keys().Source, cosigner.NewCosigner(s.cfg.DecentralizationConfig())).
		Run(context.Background(), s.cfg.MassPaymentsSenderConfig().SendingPeriod, uint64(s.cfg.MassPaymentsSenderConfig().TxsPerPeriod))
}

func (s *service) runDeferredPaymentsStreamer() {
	horizonClient := horizon.NewConnector(s.cfg.Client())
	horizonInfo, err := horizonClient.Info()
	if err != nil {
		panic(err)
	}
	builder := xdrbuild.NewBuilder(horizonInfo.Attributes.NetworkPassphrase, horizonInfo.Attributes.TxExpirationPeriod)
	processor := streamers.NewCreateDeferredPaymentRequestProcessor(s.log, pg.NewRequestsQ(s.cfg.DB()),
		horizon.NewConnector(s.cfg.Client()), 1, cosigner.NewCosigner(s.cfg.DecentralizationConfig()),
		horizonClient, builder, s.cfg.Keys().Signer, s.cfg.Keys().Source)
	dest := s.cfg.Keys().Source.Address()
	streamer := streamers.NewCreateDeferredPaymentStreamer(s.cfg.Log(), processor,
		s.cfg.Client(), streamers.RREncodeParams{
			State:        1,
			PendingTasks: s.cfg.MassPaymentsSenderConfig().Tasks,
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
