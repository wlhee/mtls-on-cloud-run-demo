# mtls-on-cloud-run-demo

This repo demos how to establish a mTLS connection to a service on Google Cloud
Run.

## Tutorial

### Build and deploy the server to Cloud Run

1. Set up env var of your GCP project ID

```
export PROJECT_ID=<your GCP project ID>

# One-time config of docker
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
  --use-http2
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




