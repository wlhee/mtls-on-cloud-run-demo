#!/bin/sh

echo "Starting TCP server ..."
/server 127.0.0.1:7777 &

echo "Waiting for TCP server to come up ..."
/nc -zv 127.0.0.1 7777
while [ $? -ne 0 ]
do
  echo "TCP server is not ready yet ..."
  sleep 0.2
  /nc -zv 127.0.0.1 7777
done
echo "TCP server is ready"

echo "Starting Envoy ..."
#/envoy -l trace -c /etc/envoy/server-envoy.yaml
/envoy -c /etc/envoy/server-envoy.yaml


