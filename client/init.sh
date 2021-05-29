#!/bin/sh

echo "Starting Envoy ..."
#/envoy -l debug -c /etc/envoy/client-envoy.yaml &
/envoy -c /etc/envoy/client-envoy.yaml &

echo "Waiting for Envoy to be ready ..."
/nc -zv 127.0.0.1 7777
while [ $? -ne 0 ]
do
  echo "Envoy is not ready yet ..."
  sleep 0.2
  /nc -zv 127.0.0.1 7777
done
echo "Envoy is ready"

echo "Starting TCP client ..."
/client 127.0.0.1:7777
