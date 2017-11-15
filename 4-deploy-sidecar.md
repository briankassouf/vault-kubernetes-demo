# Guide

## Deploy the sidecar example

This example will run two containers, the first is running consul
template that first uses the Kubernetes backend to authenticate to Vault.
It then writes secrets to a template file that is in a volume shared with the
second container. Our example app is running in the second container and loads
the template file then logs the secret (don't log secrets in a
real application!) from the template file. Consul template
will keep the secret up to date and the vault token renewed. 

### Build the sidecar example

The easiest way to build the containers is to connect your local docker agent
to the remote one in kubernetes. With minikube this can be done with:

```
eval $(minikube docker-env)
```

Then we can build the container:
```
cd sidecar/
./build_container
```

### Run the sidecar example

The `VAULT_ADDR` variable in the deployment file should be updated to your vault
server address

Now we can run the containers:

```
kubectl create -f deployment.yaml
```

### View the logs

```
kubectl logs -f  $(kubectl \
    get pods -l app=sidecar-example \
    -o jsonpath='{.items[0].metadata.name}') -c app
```

You should log output similar to:
```
2017/09/13 23:38:59 ==> WARNING: Don't ever write secrets to logs.
2017/09/13 23:38:59 ==>          This is for demonstration only.
2017/09/13 23:38:59 Username: demo
2017/09/13 23:38:59 Password: test
```

Then you should see a token renewal approximately every 20s.

### Cleanup 

When done delete the deployment and go back to the parent directory:

```
kubectl delete -f deployment.yaml
cd ..
```

## Next Steps

Read more about the kubernetes auth backend!





