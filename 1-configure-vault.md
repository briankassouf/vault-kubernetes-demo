# Deployment Guide

This tutorial will walk you configuring vault's kuberentes auth backend and
using the backend to authenticate pods.
Requirements:

* Vault 0.8.3
* Kubernetes cluster 

## Configure the Kubernetes Auth Backend 

### Configuring the backend
Mount the kubernetes auth backend:

```
vault auth-enable kubernetes
```

Configure the auth backend with the pulblic key of Kubernetes' JWT signing key,
the host for the Kubernetes API, and the CA cert used for the API. Depending on
your configuration, most of these values can be found through the `kubectl
config view` command. Replace the values below with the values for your system.

```
vault write auth/kubernetes/config \
    certificates=@/Users/brian/.minikube/apiserver.crt  \
    kubernetes_host=https://192.168.99.100:8443 \
    kubernetes_ca_cert=@/Users/brian/.minikube/ca.crt
```

### Configuring a Role

Roles are used to bind Kubernetes Service Account names and namespaces to a set
of Vault policies and token settings. 

First create the policy we want this role to gain:

```
vault policy-write kube-auth sidecar/kube-auth.hcl
```

To create a role with the S.A. name "vault-auth" in the "default" namespace:

```
vault write auth/kubernetes/role/demo \
    bound_service_account_names=vault-auth \
    bound_service_account_namespaces=default \
    policies=kube-auth \
    period=60s
```

Notice we set a period of 60s, this means the resulting token is a periodic token and
must be renewed by the application at least every 60s.

Read the demo role to verify everything was configured properly:

```
vault read auth/kubernetes/role/demo
```
Should produce the following output:
```
Key                             	Value
---                             	-----
bound_service_account_names     	[vault-auth]
bound_service_account_namespaces	[default]
max_ttl                         	0
num_uses                        	0
period                          	60
policies                        	[default kube-auth]
ttl                             	0
```

### Write a secret

This will be used later in the demo

```
vault write secret/creds username=demo password=test
```

## Next Steps

The vault server is now configured to authenticate Service Account JWT tokens
for the "vault-auth" Service Account in the "default" namespace. Next up we
will [create the Service Account in Kubernetes](./2-configure-kubernetes.md). 
