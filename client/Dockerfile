FROM envoyproxy/envoy:v1.18.2 AS envoy-base

FROM debian:8 AS debian
RUN apt-get update && apt-get -y install netcat && apt-get clean

FROM golang:1.16-buster as builder
ARG cloud_run_service_hostname

WORKDIR /
COPY client.go .
RUN go build -v -o client client.go

COPY --from=envoy-base /usr/local/bin/envoy /envoy
COPY client-envoy.yaml /etc/envoy/client-envoy.yaml
RUN sed -i "s|CLOUD_RUN_SERVICE_HOSTNAME|${cloud_run_service_hostname}|g" /etc/envoy/client-envoy.yaml

COPY certs/cert.pem /etc/certs/cert.pem
COPY certs/key.pem /etc/certs/key.pem
COPY certs/ca-certificates.crt /etc/certs/ca-certificates.crt

COPY init.sh .

COPY --from=debian /bin/nc nc

CMD ["/init.sh"]