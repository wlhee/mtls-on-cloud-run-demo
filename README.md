# mtls-on-cloud-run-demo

This repo demos how to establish a mTLS connection to a service on Google Cloud
Run.

## Tutorial

### Build and deploy the server to Cloud Run

1. Set up env var of your GCP project ID

```
export PROJECT_ID=<your GCP project ID>

# One-time config for docker push to Google Container Reigstry.
gcloud auth configure-docker

```

Build the server

```
cd server/

docker build -t gcr.io/$PROJECT_ID/mtls-demo .

docker push gcr.io/$PROJECT_ID/mtls-demo
```

Deploy to Cloud Run

```
gcloud beta run deploy mtls-demo \
  --image gcr.io/$PROJECT_ID/mtls-demo \
  --allow-unauthenticated \
  --use-http2 \
  --project ${PROJECT_ID}
```

Remember the URL of the newly deployed service
`https://mtls-demo-<hash>-uc.a.run.app`

### 2. Build and connect the client with mTLS to the server

Build the client with the *hostname* (excluding the `https://` scheme) of the
Cloud Run service depployed above.

```
cd client/

export SERVICE_HOSTNAME=<your Cloud Run service hostname without scheme>

docker build \
--build-arg cloud_run_service_hostname=${SERVICE_HOSTNAME} \
-t gcr.io/$PROJECT_ID/client .
```

Run the clinet

```
docker run --network=host gcr.io/$PROJECT_ID/client
```

Check the result

```
Starting the client ...
================================================================
== Congrats!                                                  ==
== If you see this message, it means                          ==
== you've successfully run the mTLS demo on Google Cloud Run! ==
================================================================
```

## How it works

This demo leverages the advanced capability of envoy proxy to
[tunnel TCP over HTTP POST](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/http/upgrades#tunneling-tcp-over-http).

The following diagram illustrates the how the tunnel works:

```
TCP client                                                TCP server
   |                                                         |
   |-> raw                                             raw ->|
       TCP ->|                                     |-> TCP
             |                                     |
         client-sdie                           server-side
           envoy                                 envoy
             |                                     |
             |-> mTLS                       mTLS ->|
                 TCP ->|                |-> TCP
                       |                |
                  client-sdie       server-side
                     envoy             envoy
                       |                |
                       |----> TLS ----->|
                           HTTP/2 POST
```

In the demo above, the server image contains the server-side envoy and the TCP
server. The image is deployed to Cloud Run with HTTP/2 enabled.

The cient image contains the TCP client and the client-side envoy. It runs
locally and connects to the service running on Cloud Run with a TCP tunnel
over HTTP/2 POST stream.
