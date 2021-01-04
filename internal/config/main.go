package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/tokend/connectors/builder"
	"gitlab.com/tokend/connectors/keyer"
	"gitlab.com/tokend/connectors/signed"
	"gitlab.com/tokend/connectors/submit"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	signed.Clienter
	keyer.Keyer
	submit.Submission
	builder.Builderer
	Doorman
	MassPaymentsSenderConfiger
	Decentralizationer
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	types.Copuser
	comfig.Listenerer
	signed.Clienter
	keyer.Keyer
	submit.Submission
	builder.Builderer
	Doorman
	MassPaymentsSenderConfiger
	Decentralizationer
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:                     getter,
		Databaser:                  pgdb.NewDatabaser(getter),
		Copuser:                    copus.NewCopuser(getter),
		Listenerer:                 comfig.NewListenerer(getter),
		Logger:                     comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Clienter:                   signed.NewClienter(getter),
		Keyer:                      keyer.NewKeyer(getter),
		Doorman:                    NewDoorman(getter),
		Submission:                 submit.NewSubmission(getter),
		Builderer:                  builder.NewBuilderer(getter),
		MassPaymentsSenderConfiger: NewMassPaymentsSenderConfiger(getter),
		Decentralizationer:         NewDecentralizationer(getter),
	}
}
