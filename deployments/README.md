# Deploying to Google Kubernetes Engine (GKE Autopilot)
This guide describes how to deploy ChatOrbit to Google Kubernetes Engine (GKE) in Autopilot mode.

## What's included
- Namespace isolation via `chatorbit-prod`.
- Stateful workloads for MongoDB and Redis with persistent volumes.
- Deployments and ClusterIP services for the `auth`, `chat`, and `user` microservices.
- Centralized configuration and secrets via `ConfigMap` and `Secret`.
- A `kustomization.yaml` that applies everything in one command.


### Enable the Kubernetes Engine API
1. Open the Google Cloud Console.
2. Search for **"Kubernetes Engine"**.
3. Enable the API if prompted.

### Create a GKE Autopilot Cluster
1. In the Google Cloud Console, go to **Kubernetes Engine**.
2. Click **Create**.
3. Choose **Autopilot**.
4. Choose region: `asia-east1`.

**Why Autopilot mode?**
You do not need to configure nodes, node pools, instance types, autoscaling rules, or VPC details—Google handles all of it. Autopilot automatically manages nodes, autoscaling, networking, and compute.

### Install Google Cloud SDK on Ubuntu
```bash
sudo snap install google-cloud-sdk --classic
```
Check installation:
```bash
gcloud --version
```
Log in:
```bash
gcloud auth login
```

### Connect kubectl
Fetch credentials for your cluster:
```bash
gcloud container clusters get-credentials CLUSTER_NAME --region asia-east1
```
Verify connectivity:
```bash
kubectl get nodes
```

If you see the error `gke-gcloud-auth-plugin, which is needed for continued use of kubectl, was not found or is not executable`, install the plugin before connecting:
```bash
sudo snap install google-cloud-cli
sudo apt-get install google-cloud-sdk-gke-gcloud-auth-plugin
```
Then retry:
```bash
kubectl get nodes
kubectl get pods
```
You should see the Autopilot node pool listed.

### Deploy the app

Apply the Kubernetes manifests in the `deployments/prod/k8s` directory:
```bash
kubectl apply -f deployments/prod/k8s
```
Check status:
```bash
kubectl get pods
```

After applying the manifests, set your production JWT secret:
```bash
kubectl create ns chatorbit-prod --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f deployments/prod/k8s/secret.yaml --namespace chatorbit-prod
kubectl patch secret chatorbit-secrets -n chatorbit-prod \
  --type merge -p '{"stringData": {"JWT_SECRET": "<strong-secret-value>"}}'
```
Verify:
```bash
kubectl get secret chatorbit-secrets -n chatorbit-prod -o yaml
```
### Get External Public URLs
Each backend service exposes a LoadBalancer.
```bash
kubectl get svc -n chatorbit-prod
```
Use these external IPs in the frontend.

### Local Debugging
Test locally without LoadBalancer:
```bash
kubectl port-forward -n chatorbit-prod svc/auth-service 8089:8089
kubectl port-forward -n chatorbit-prod svc/user-service 8087:8087
kubectl port-forward -n chatorbit-prod svc/chat-service 8088:8088
```
