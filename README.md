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
3) Go to `System --> Kubernetes Cluster` and create a new one
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
apoctl api import --file apowine/policy/policy.yaml -n /$APORETO_ACCOUNT/$K8S_CLUSTER_NAME/apowine
```

If running against pre-production, specify `policy-alpha.yaml`.


### Network Layer

There are 4 network policies defined in the `policy.yaml` file:
1) Allow anyone to access the UI (`allow-access-from-external`)
2) Allow the UI to access midgard
By default, this is going to midgard that is running on production.
This means the `JWTSigningCertificate` defined in the service, should point to the midgard that is available in production.
3) Allow anyone to do DNS resolution (`allow-dns-resolutions`)
4) Any PU can talk to `mongodb` server or `apowine-server` service.


_Note: we should provide a way to modify the midgard url using `APOWINE_MIDGARD_URL`_

### Service Layer

There are 3 services defined in the `policy.yaml` file:
1) `apowine-mongodb` is exposing the database on port 27017
2) `apowine-server` is exposed internally to the cluster at `server.apowine.svc.cluster.local:3000`.
It has a `JWTSigningCertificate` to be able to verify the token and decode it to retrieve its claims.
3) `apowine-ui` is publicly exposed on `apowine.aporeto.com` on port `4443`.


### How to access the UI ?

Find the External IP of the `client-public` service that is exposed via a Load Balancer:
```
kubectl -n apowine get svc
```

Edit your `/etc/hosts` to add an entry to apowine.aporeto.com:
```
sudo vim /etc/hosts
...
35.232.92.93    apowine.aporeto.com
```

Use your browser to access the UI on https://apowine.aporeto.com:4443/.
