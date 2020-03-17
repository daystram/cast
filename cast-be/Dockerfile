FROM daystram/go-builder:1.13 AS builder
WORKDIR $GOPATH/src/gitlab.com/daystram/cast/cast-be
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /cast-be .

FROM alpine:3.10
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /cast-be /app/
ENTRYPOINT ["/app/cast-be"]
