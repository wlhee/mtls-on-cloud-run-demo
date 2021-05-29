FROM envoyproxy/envoy:v1.18.2 AS envoy-base

FROM debian:8 AS debian
RUN apt-get update && apt-get -y install netcat && apt-get clean

FROM golang:1.16-buster as builder

WORKDIR /
COPY server.go .
RUN go build -v -o server server.go

COPY --from=envoy-base /usr/local/bin/envoy /envoy
COPY server-envoy.yaml /etc/envoy/server-envoy.yaml
COPY certs/cert.pem /etc/certs/cert.pem
COPY certs/key.pem /etc/certs/key.pem
EXPOSE 8080

COPY init.sh .

COPY --from=debian /bin/nc nc

CMD ["/init.sh"]