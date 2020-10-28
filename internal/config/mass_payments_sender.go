package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type MassPaymentsSenderConfig struct {
	SendingPeriod int64  `figure:"sending_period"`
	TxsPerPeriod  int64  `figure:"txs_per_period"`
	Tasks         uint32 `figure:"tasks"`
}

type MassPaymentsSenderConfiger interface {
	MassPaymentsSenderConfig() MassPaymentsSenderConfig
}

func NewMassPaymentsSenderConfiger(getter kv.Getter) MassPaymentsSenderConfiger {
	return &massPaymentsSender{
		getter: getter,
	}
}

type massPaymentsSender struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *massPaymentsSender) MassPaymentsSenderConfig() MassPaymentsSenderConfig {
	return c.once.Do(func() interface{} {
		config := MassPaymentsSenderConfig{
			SendingPeriod: 5,
			TxsPerPeriod:  5,
			Tasks:         1,
		}

		err := figure.
			Out(&config).
			From(kv.MustGetStringMap(c.getter, "mass_payments_sender")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out client"))
		}

		return config
	}).(MassPaymentsSenderConfig)
}
