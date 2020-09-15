package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Doorman interface {
	SkipSignCheck() bool
}

func NewDoorman(getter kv.Getter) Doorman {
	return &doorman{
		getter: getter,
	}
}

type doorman struct {
	getter kv.Getter
	once   comfig.Once
}

func (d *doorman) SkipSignCheck() bool {
	return d.once.Do(func() interface{} {
		config := struct {
			SkipSignCheck bool `fig:"skip_sign_check"`
		}{
			SkipSignCheck: false,
		}

		raw, err := d.getter.GetStringMap("doorman")
		if err != nil {
			raw = make(map[string]interface{})
		}
		err = figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return config.SkipSignCheck
	}).(bool)
}
