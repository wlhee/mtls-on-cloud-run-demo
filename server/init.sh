#!/bin/sh

echo "Starting the server ..."
/server 127.0.0.1:7777 &

echo "Waiting for the server to come up ..."
/nc -zv 127.0.0.1 7777
while [ $? -ne 0 ]
do
  sleep 0.2
  /nc -zv 127.0.0.1 7777
done

echo "Starting Envoy ..."
#/envoy -l debug -c /etc/envoy/server-envoy.yaml
/envoy -c /etc/envoy/server-envoy.yaml


