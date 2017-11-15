# Guide

This guide will walk you configuring Vault's Kubernetes Auth Backend and
using the backend to authenticate applications.
Requirements:

* Vault >0.9.0
* Kubernetes cluster 

## Configure a Kubernetes Service Account for Verifying JWTs

The Kubernetes Authentication Backend has a `token_reviewer_jwt` field which 
takes a Service Account Token that is in charge of verifying the validity of 
other service account tokens. This Service Account will call into Kubernetes
TokenReview API and verify the service account tokens provided during login
are valid. 

### Prerequisites

Vault uses the Kubernetes TokenReview API to validate that JWT tokens are still
valid, and have not been deleted from within Kubernetes.

To ensure stale/deleted Service Accounts tokens can not authenticate with vault
the Kubernetes API server must be running with `--service-account-lookup`. This
is defaulted to on in Kubernetes 1.7 but prior versions should ensure this is
set.

### Create the Service Account

The service account is defined in `vault-reviewer.yaml` and can be created with this
command:

```
kubectl create -f vault-reviewer.yaml
```

### RBAC 

If your kubernetes cluster uses RBAC authorization you will need to provide the
service account with a role that gives it access to the TokenReview API. If not, 
this step can be skipped.

The RBAC role is defined in `vault-reviewer-rbac.yaml` and can be created with
this command:

```
kubectl create -f vault-reviewer-rbac.yaml
```

## Configure a Kubernetes Service Account for login in

This service account will be used to login to the auth backend.

### Create the Service Account

The service account is defined in `vault-auth.yaml` and can be created with this
command:

```
kubectl create -f vault-auth.yaml
```

This service account does not need any RBAC permissions.

## Next Steps

We now have a service account setup with the appropriate permissions. Next we
will [configure the Kubernetes Auth Backend](./2-configure-vault.md).
