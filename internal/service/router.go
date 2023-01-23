package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/doorman"
	"github.com/SafeRE-IT/mass-payments-sender-svc/internal/data/pg"
	"github.com/SafeRE-IT/mass-payments-sender-svc/internal/horizon"
	"github.com/SafeRE-IT/mass-payments-sender-svc/internal/service/handlers"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()
	horizonClient := horizon.NewConnector(s.cfg.Client())
	state, err := horizonClient.Info()
	if err != nil {
		panic(errors.Wrap(err, "failed to get horizon info"))
	}

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxRequestsQ(pg.NewRequestsQ(s.cfg.DB())),
			handlers.CtxTransactionsQ(pg.NewPaymentsQ(s.cfg.DB())),
			handlers.CtxHorizonInfo(state),
			handlers.CtxDoorman(doorman.New(
				s.cfg.SkipSignCheck(),
				horizonClient),
			),
			handlers.CtxKeys(s.cfg.Keys().Source),
			handlers.CtxSigner(s.cfg.Keys().Signer),
			handlers.CtxXdrBuilder(s.cfg.Builder()),
		),
	)

	r.Route("/integrations/mass-payments", func(r chi.Router) {
		r.Get("/info", handlers.GetInfo)
		r.Route("/requests", func(r chi.Router) {
			r.Get("/", handlers.GetRequestsList)
			r.Get("/{id}", handlers.GetRequest)
		})
		r.Route("/payments", func(r chi.Router) {
			r.Get("/", handlers.GetPaymentsList)
			r.Get("/csv", handlers.GetPaymentsListCsv)
		})
		r.Post("/sign", handlers.SignTx)
	})

	return r
}
