FROM golang as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY ./ ./
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /bin/centre /app/cmd/centre.go

FROM scratch as prod

ENV PROBE_DSN="host=postgres user=test dbname=test port=5432"
ENV GIN_MODE="debug"
COPY --from=0 /bin/centre .
ENTRYPOINT ["/centre"]

EXPOSE 12000