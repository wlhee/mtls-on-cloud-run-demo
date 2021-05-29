#!/bin/sh

echo "Starting Envoy ..."
/envoy -l debug -c /etc/envoy/server-envoy.yaml &
#/envoy -c /etc/envoy/server-envoy.yaml &


echo "Starting the server ..."
/server

