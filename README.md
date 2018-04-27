# apowine

Apowine is a simple application that demonstrates how aporeto uses policies to secure  API endpoints between containers in a kubernetes cluster.


There are 3 major and an associated components in apowine

----
**More on Kubernetes network policies:**

* Frontend-UI/Frontend-Curl

      Frontend-UI is a simple UI that queries the server

      Frontend-Curl will use /GET endpoint to request beers/wines every 30 seconds

      Both are deployed as a separate containers and based on your preference you can deploy either

* Backend/Server

      Server handles the necessary backend logic for the client to manipulate data in DB

* Database

      MongoDB is used as the database and uses default credentials


* Producer-Beer/Producer-Wine

      Producer basically prepopulates the database with beers/wines by calling the server endpoint


## Steps to deploy

 Apowine is easy to deploy on kubernetes without any parameters or environment variables

 Each and every component should communicate by default

1) Checkout the deployment files:
```
git clone https://github.com/aporeto-inc/apowine.git
cd apowine/deployments
```

2) Deploy all the other components that uses this namespace
```
kubectl create -f .
```
