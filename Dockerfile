FROM golang as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY ./ ./
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /probe /app/cmd/centre.go

FROM scratch as prod

ENV PROBE_DSN="host=postgres user=test dbname=test port=5432"
ENV GIN_MODE="debug"
COPY --from=0 /probe .
ENTRYPOINT ["/probe"]

EXPOSE 12000