FROM golang:alpine AS builder

WORKDIR $GOPATH/src/dreddick.home/es_client/
COPY main.go .

# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 go build


FROM scratch
# Copy our static executable.
COPY --from=builder /go/src/dreddick.home/es_client/es_client /go/bin/es_client
COPY  eu-west-1-es-amazonaws-com-chain.pem /etc/ssl/certs/
ENTRYPOINT ["/go/bin/es_client"]
