package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type DecentralizationConfig struct {
	Disabled      bool     `fig:"disabled"`
	IsCoordinator bool     `fig:"is_coordinator"`
	Nodes         []string `fig:"nodes"`
}

type Decentralizationer interface {
	DecentralizationConfig() DecentralizationConfig
}

func NewDecentralizationer(getter kv.Getter) Decentralizationer {
	return &decentralizationer{
		getter: getter,
	}
}

type decentralizationer struct {
	getter kv.Getter
	once   comfig.Once
}

func (d *decentralizationer) DecentralizationConfig() DecentralizationConfig {
	return d.once.Do(func() interface{} {
		config := DecentralizationConfig{}

		raw, err := d.getter.GetStringMap("decentralization")
		if err != nil {
			config.Disabled = true
			return config
		}
		err = figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}

		return config
	}).(DecentralizationConfig)
}
