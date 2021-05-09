FROM scratch

COPY ./probe /probe
ENTRYPOINT ["/probe"]
ENV PROBE_DSN="host=192.168.1.119 user=probe dbname=probe port=25432"
ENV GIN_MODE="debug"

EXPOSE 12000