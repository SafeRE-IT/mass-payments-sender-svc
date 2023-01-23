package handlers

import (
	"context"
	"net/http"

	"gitlab.com/tokend/go/xdrbuild"

	"gitlab.com/tokend/keypair"

	"gitlab.com/tokend/go/doorman"
	"github.com/SafeRE-IT/mass-payments-sender-svc/internal/data"
	regources "gitlab.com/tokend/regources/generated"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	requestsQCtxKey
	paymentsQCtxKey
	horizonStateCtxKey
	doormanCtxKey
	keysCtxKey
	signerCtxKey
	xdrBuilderCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxRequestsQ(entry data.RequestsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestsQCtxKey, entry)
	}
}

func RequestsQ(r *http.Request) data.RequestsQ {
	return r.Context().Value(requestsQCtxKey).(data.RequestsQ).New()
}

func CtxTransactionsQ(entry data.PaymentsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, paymentsQCtxKey, entry)
	}
}

func TransactionsQ(r *http.Request) data.PaymentsQ {
	return r.Context().Value(paymentsQCtxKey).(data.PaymentsQ).New()
}

func CtxHorizonInfo(entry regources.HorizonState) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, horizonStateCtxKey, entry)
	}
}

func HorizonInfo(r *http.Request) regources.HorizonState {
	return r.Context().Value(horizonStateCtxKey).(regources.HorizonState)
}

func CtxDoorman(d doorman.Doorman) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, doormanCtxKey, d)
	}
}

func Doorman(r *http.Request, constraints ...doorman.SignerConstraint) error {
	d := r.Context().Value(doormanCtxKey).(doorman.Doorman)
	return d.Check(r, constraints...)
}

func CtxKeys(entry keypair.Address) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, keysCtxKey, entry)
	}
}

func Keys(r *http.Request) keypair.Address {
	return r.Context().Value(keysCtxKey).(keypair.Address)
}

func CtxSigner(entry keypair.Full) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, signerCtxKey, entry)
	}
}

func Signer(r *http.Request) keypair.Full {
	return r.Context().Value(signerCtxKey).(keypair.Full)
}

func CtxXdrBuilder(entry *xdrbuild.Builder) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, xdrBuilderCtxKey, entry)
	}
}

func XdrBuilder(r *http.Request) *xdrbuild.Builder {
	return r.Context().Value(xdrBuilderCtxKey).(*xdrbuild.Builder)
}
