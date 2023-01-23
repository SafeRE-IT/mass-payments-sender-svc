FROM golang:1.18-alpine

RUN apk add --no-cache git build-base

WORKDIR /go/src/github.com/SafeRE-IT/mass-payments-sender-svc
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/mass-payments-sender-svc github.com/SafeRE-IT/mass-payments-sender-svc

###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/mass-payments-sender-svc /usr/local/bin/mass-payments-sender-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["mass-payments-sender-svc"]
