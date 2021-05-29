FROM envoyproxy/envoy:v1.18.2 AS envoy-base

FROM golang:1.16-buster as builder

WORKDIR /
COPY client.go .
RUN go build -v -o client client.go

COPY --from=envoy-base /usr/local/bin/envoy /envoy
COPY client-envoy.yaml /etc/envoy/client-envoy.yaml
COPY cert.pem /etc/certs/cert.pem
COPY key.pem /etc/certs/key.pem
COPY ca-certificates.crt /etc/certs/ca-certificates.crt

COPY init.sh .

CMD ["/init.sh"]