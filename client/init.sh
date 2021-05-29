#!/bin/sh

echo "Starting Envoy ..."
#/envoy -l debug -c /etc/envoy/client-envoy.yaml &
/envoy -c /etc/envoy/client-envoy.yaml &

echo "Waiting for Envoy to be ready ..."
sleep 2

echo "Starting the client ..."
/client 127.0.0.1:7777
