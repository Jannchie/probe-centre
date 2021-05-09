FROM golang as builder

WORKDIR /app

COPY ./ ./
RUN GOOS=linux GOARCH=amd64 go build -o /probe /app/cmd/centre.go

FROM scratch as prod

ENV PROBE_DSN="host=postgres user=test dbname=test port=5432"
ENV GIN_MODE="debug"
COPY --from=0 /probe .
ENTRYPOINT ["/probe"]

EXPOSE 12000