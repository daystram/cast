FROM daystram/go-builder:1.13 AS builder
WORKDIR $GOPATH/src/gitlab.com/daystram/cast/cast-is
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /cast-is .

FROM daystram/ingest-base:latest
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /cast-is /app/
COPY config/app.conf /app/
ENTRYPOINT ["/app/cast-is"]
