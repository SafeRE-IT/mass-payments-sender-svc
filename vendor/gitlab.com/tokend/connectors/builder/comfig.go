package builder

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/tokend/connectors/lazyinfo"
	"gitlab.com/tokend/connectors/signed"
	"gitlab.com/tokend/go/xdrbuild"
)

type Builderer interface {
	Builder() *xdrbuild.Builder
}

type builderer struct {
	signed.Clienter

	getter kv.Getter
	once   comfig.Once
}

func NewBuilderer(getter kv.Getter) *builderer {
	return &builderer{
		getter:   getter,
		Clienter: signed.NewClienter(getter),
	}
}

func (h *builderer) Builder() *xdrbuild.Builder {
	return h.once.Do(func() interface{} {

		cli := h.Clienter.Client()
		info, err := lazyinfo.New(cli).Info()
		if err != nil {
			panic(errors.Wrap(err, "failed to get info"))
		}

		return xdrbuild.NewBuilder(info.Attributes.NetworkPassphrase, info.Attributes.TxExpirationPeriod)
	}).(*xdrbuild.Builder)
}
