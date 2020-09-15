package horizon

import (
	jsonapi "gitlab.com/distributed_lab/json-api-connector"
	"gitlab.com/distributed_lab/json-api-connector/client"
	"gitlab.com/tokend/connectors/keyvalue"
	"gitlab.com/tokend/connectors/lazyinfo"
	"gitlab.com/tokend/connectors/submit"
)

type Connector struct {
	connector *jsonapi.Connector

	*lazyinfo.LazyInfoer
	*keyvalue.KeyValuer
	*submit.Submitter
}

func NewConnector(client client.Client) *Connector {
	return &Connector{
		connector:  jsonapi.NewConnector(client),
		LazyInfoer: lazyinfo.New(client),
		KeyValuer:  keyvalue.New(client),
		Submitter:  submit.New(client),
	}
}
