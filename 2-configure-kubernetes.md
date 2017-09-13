# Deployment Guide

## Configure the Kubernetes Service Account

### Prerequisites
Vault uses the Kubernetes TokenReview API to validate that JWT tokens are still
valid, and have not been deleted from within Kubernetes.

To ensure stale/deleted Service Accounts tokens can not authenticate with vault
the Kubernetes API server must be running with `--service-account-lookup`. This
is defaulted to on in Kubernetes 1.7 but prior versions should ensure this is
set.

### Create the Service Account

The service account is defined in `vault-auth.yaml` and can be created with this
command:

```
kubectl create -f vault-auth.yaml
```

### RBAC 

If your kubernetes cluster uses RBAC authorization you will need to provide the
service account with a role that gives it access to the TokenReview API. 

The RBAC role is defined in in `vault-auth-rbac.yaml` and can be created with
this command:

```
kubectl create -f vault-auth-rbac.yaml
```

## Next Steps

We now have a configure vault backend and the service account setup with the
appropriate permissions. Next we can deploy a basic application.
