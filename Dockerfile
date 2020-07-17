# Use a mutli-stage build pipeline to generate the executable
FROM golang:1.14.6

ARG VERSION="development"

ENV GO_PATH="/go"

ADD . $GO_PATH/src/github.com/SierraSoftworks/minback-cleanup
WORKDIR $GO_PATH/src/github.com/SierraSoftworks/minback-cleanup

RUN go get -t ./...
RUN go test -v ./...

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o bin/minback-cleanup -a -installsuffix cgo -ldflags "-s -X main.version=$VERSION"

# Build the actual container
FROM alpine:latest
LABEL maintainer="Benjamin Pannell <admin@sierrasoftworks.com>"

RUN apk add --update tini
ENTRYPOINT ["/sbin/tini", "--"]

COPY --from=0 /go/src/github.com/SierraSoftworks/minback-cleanup/bin/minback-cleanup /bin/minback-cleanup

LABEL VERSION=$VERSION

WORKDIR /bin

ENV MINIO_SERVER=""
ENV MINIO_BUCKET="backups"
ENV MINIO_ACCESS_KEY=""
ENV MINIO_SECRET_KEY=""

ENTRYPOINT ["/bin/minback-cleanup"]
CMD ["cleanup"]