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

Remember the URL of the newly deployed service, for example:
`https://mtls-demo-<hash>-uc.a.run.app`

### 2. Build and connect the client with mTLS to the server

Build the client with the *hostname* (excluding the `https://` scheme) of the
Cloud Run service depployed above.

```
cd client/

# Example: export SERVICE_HOSTNAME=mtls-demo-<hash>-uc.a.run.app
export SERVICE_HOSTNAME=<your Cloud Run service hostname without scheme>

docker build \
--build-arg cloud_run_service_hostname=${SERVICE_HOSTNAME} \
-t gcr.io/$PROJECT_ID/client .
```

Run the client

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
         client-side                           server-side
           envoy                                 envoy
             |                                     |
             |-> mTLS                       mTLS ->|
                 TCP ->|                |-> TCP
                       |                |
                  client-side       server-side
                     envoy             envoy
                       |                |
                       |----> TLS ----->|
                           HTTP/2 POST
```

In the demo above, the server image contains the server-side envoy and the TCP
server. The image is deployed to Cloud Run with HTTP/2 enabled.

The client image contains the TCP client and the client-side envoy. It runs
locally and connects to the service running on Cloud Run with a TCP tunnel
over HTTP/2 POST stream.

## Implication for load balancing on Cloud Run

Cloud Run load balances each incoming HTTP/2 streams (and HTTP/1 requests)
across the serving instances. In this demo, as each TCP stream is tunneled
over a HTTP/2 stream, from Cloud Run point of view, there is **only one**
"active request", which is served by a single instance. If the client
application multiplexes multiple "requests" into one TCP stream, they will
be sent to the same instance instead of load balanced across multiple
instances. The client is responsible for deciding the level of mutiplexing
over one TCP stream.

For example, if the client is a HTTP/1 client, as HTTP/1 doesn't support
multiplexing, each HTTP/1 request is tunneled over a separate TCP connection,
which in turns is tunneled over a separate HTTP/2 stream. Cloud Run can still
load balance the HTTP/1 requests across multiple instances.

If the client is a HTTP/2 client, as HTTP/2 supports multiplexing multiple
HTTP/2 streams into a same TCP connection, multiple HTTP/2 streams will be
tunneled over a single HTTP/2 stream, which is serverd by one Cloud Run
instance. To mitigate this issue, the HTTP/2 client can lower than the
multiplexing level, i.e, configuring max streams per connection to a smaller
value.

## What if Cloud Run service has auth enabled?

Cloud Run supports "private service", which requires authentication and
authorization for accessing the serivce. The following shows how to make the
client work with a auth-enabled Cloud Run service.

1. Deploy the server to Cloud Run with auth enabled

```
cd server/

gcloud beta run deploy mtls-demo-auth \
  --image gcr.io/$PROJECT_ID/mtls-demo \
  --no-allow-unauthenticated \
  --use-http2 \
  --project ${PROJECT_ID}
```

Remember the URL of the newly deployed service, for example:
`https://mtls-demo-auth-<hash>-uc.a.run.app`

Test the service with a simple request, 403 is expected:

```
curl -i https://mtls-demo-auth-<hash>-uc.a.run.app

HTTP/2 403 
...
```

Test the service with a request with identity token header, 404 is expected.
If not, follow this [doc](https://cloud.google.com/run/docs/authenticating/developers)
to set up proper auth.

```
curl -i -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
https://mtls-demo-auth-<hash>-uc.a.run.app

HTTP/2 404
...
```

2. Run a simple auth server that can provide ID token to the client

```
cd client/auth_server

go run server.go

Authz server is listening on port 8080
```

3. Build a new client that fetches an ID token from the auth server. The new
client uses the envoy [ext_authz fitler](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter)
that makes an HTTP request to the auth server for fetching the ID token
and inject it as authorization header before sending the actual request to
Cloud Run.

```
cd client/

# Example: export SERVICE_HOSTNAME=mtls-demo-auth-<hash>-uc.a.run.app
export SERVICE_HOSTNAME=<your Cloud Run service hostname without scheme>

docker build \
--build-arg cloud_run_service_hostname=${SERVICE_HOSTNAME} \
-t gcr.io/$PROJECT_ID/client-auth \
-f Dockerfile_auth .
```

4. Run the client and check the result

```
docker run --network=host gcr.io/$PROJECT_ID/client-auth

Starting the client ...
================================================================
== Congrats!                                                  ==
== If you see this message, it means                          ==
== you've successfully run the mTLS demo on Google Cloud Run! ==
================================================================
```




