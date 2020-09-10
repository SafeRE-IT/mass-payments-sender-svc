FROM golang:1.12

WORKDIR /go/src/gitlab.com/tokend/mass-payments-sender-svc
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/mass-payments-sender-svc gitlab.com/tokend/mass-payments-sender-svc

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/mass-payments-sender-svc /usr/local/bin/mass-payments-sender-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["mass-payments-sender-svc"]
