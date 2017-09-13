# Deployment Guide

## Deploy the basic example

This example application will read the servcie account JWT token and use it to
authenticate with Vault. It will then log the token (don't log secrets in a
normal use case!) and keep the token renewed.

### Build the basic example

The easiest way to build the container is to connect your local docker agent
to the remote one in kubernetes. With minikube this can be done with:

```
eval $(minikube docker-env)`
```

Then we can build the container:
```
cd basic/
./build_container
```

### Run the basic example

Now we can run the application:

```
kubectl create -f deployment.yaml
```

### View the logs

```
kubectl logs -f  $(kubectl \
    get pods -l app=basic-example \
    -o jsonpath='{.items[0].metadata.name}')
```

You should log output similar to:
```
2017/09/13 23:17:06 ==> WARNING: Don't ever write secrets to logs.
2017/09/13 23:17:06 ==>          This is for demonstration only.
2017/09/13 23:17:06 b18c700d-2577-8b08-d0c7-49e01aefd8f3
2017/09/13 23:17:06 Starting renewal loop
```

Then you should see a token renewal approximately every 20s.

### Cleanup 

When done delete the deployment and go back to the parent directory:

```
kubectl delete -f deployment.yaml
cd ..
```

## Next Steps

Next we will run a sidecar example.




