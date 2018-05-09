# apowine

Apowine is a simple application that demonstrates how aporeto uses policies to secure
apis between containers in a kubernetes cluster using Aporeto Services APIs.

## Components

There are 3 majors and an associated components in apowine

### Frontend-UI/Frontend-Curl

- Frontend-UI is a simple UI that queries the server
- Frontend-Curl will use /GET endpoint to request beers/wines every 30 seconds
- Server handles the necessary backend logic for the client to manipulate data in DB.
- MongoDB is used as the database and uses default credentials.
- Producers basically prepopulates the database with beers/wines by calling the server endpoint

## Deployment

Apowine is easy to deploy on kubernetes without any parameters or environment variables

### Aporeto deployment

> If you already have one k8s cluster running Aporeto enforcerd you can skip this part.

1) Create a kubernetes cluster
2) Log into your Aporeto account
3) Go to `Kubernetes Cluster` and create a new one
4) Unzip the downloaded tgz
5) run `kubectl apply -f /path/to/downloaded/folder`

If your cluster runs on GCP and have permission issue while, you must run

```bash
kubectl create clusterrolebinding cluster-admin-binding --clusterrole=cluster-admin --user=$GCP_USER_EMAIL`
```

Wait a few seconds for all instances of enforcerd to be ready.

### Apowine deployment

First clone the repository:

```bash
git clone https://github.com/aporeto-inc/apowine.git
```

Then deploy apowine:

```bash
kubectl apply -f apowine/deployments/kubernetes
```

Then install the aporeto policy

```bash
apoctl import --file apowine/policy/policy.yml -n /$APORETO_ACCOUNT/$K8S_CLUSTER_NAME/apowine --mode full
```
