FROM golang:alpine AS builder

WORKDIR $GOPATH/src/devopsgoat/go-elasticsearch-duster/
COPY main.go .

# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 go build


FROM scratch
# Copy our static executable.
COPY --from=builder /go/src/devopsgoat/go-elasticsearch-duster/go-elasticsearch-duster /go/bin/app
COPY  eu-west-1-es-amazonaws-com-chain.pem /etc/ssl/certs/
ENTRYPOINT ["/go/bin/app"]
