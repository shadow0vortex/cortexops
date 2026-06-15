# CortexOps GKE Autopilot Deployment Blueprint

This document outlines the architecture, requirements, and step-by-step instructions for deploying CortexOps to Google Kubernetes Engine (GKE) Autopilot for the `cortexops.amshithnair.in` domain.

## 1. Exact GCP Services Required
- **GKE Autopilot**: Fully managed Kubernetes cluster for workloads.
- **Artifact Registry**: Private Docker image storage for CortexOps microservices.
- **Cloud SQL for PostgreSQL (Optional but Recommended)**: Managed Postgres 15 for production instead of in-cluster Postgres.
- **VPC Network**: Custom VPC, subnets, and Cloud NAT for outbound internet access from private nodes.
- **Cloud Load Balancing**: Automatically provisioned by GKE Ingress controller for exposing services.

## 2. Estimated Monthly Cost (GCP + Cloudflare)
- **GKE Autopilot**: Pay-per-pod resource usage (~$50 - $100/mo depending on traffic, base fee of $73/mo is usually covered by free tier for first cluster).
- **Artifact Registry**: ~$1-5/mo (storage & network egress).
- **Cloud SQL (PostgreSQL)**: ~$50/mo (db-custom-1-3840 for a small production setup).
- **Persistent Disks**: ~$10-20/mo (Qdrant & Temporal storage).
- **Cloudflare**: $0 (Free tier) or $20/mo (Pro tier).
- **Total Estimated**: $110 - $200 / month.

## 3. Cluster Sizing Recommendations
Since GKE Autopilot automatically scales nodes based on pod specifications, we focus on pod sizing. Based on `values.yaml`:
- **Collector/Correlator/Remediation**: `requests: 200m CPU, 256Mi RAM`, `limits: 500m CPU, 512Mi RAM`.
- **Topology**: `requests: 100m CPU, 512Mi RAM`, `limits: 200m CPU, 1Gi RAM`.
- **RCA (AI workloads)**: `requests: 500m CPU, 1Gi RAM`, `limits: 1000m CPU, 2Gi RAM`.
- **Autoscaling**: Configure Horizontal Pod Autoscaler (HPA) to scale between 2 and 10 replicas per microservice based on 80% CPU/Memory utilization.

## 4. Required Kubernetes Resources
- **Deployments**: Collector, Correlator, Topology, RCA, Remediation, Temporal workers.
- **StatefulSets**: Qdrant (Vector DB), NATS JetStream, in-cluster PostgreSQL (if not using Cloud SQL).
- **Services**: ClusterIP for internal communication.
- **Ingress**: GKE Ingress (Google Cloud Load Balancer) mapping `api.*`, `grafana.*`, `temporal.*` to respective services.
- **ServiceAccounts**: `cortexops-sa` with restricted privileges.
- **Secrets**: DB credentials, NATS auth, Temporal TLS certs, Cloudflare API tokens (for cert-manager).

## 5. Required Persistent Volumes (PVs)
- **PostgreSQL**: Standard SSD (pd-ssd), 50GB. (If self-hosted in cluster).
- **Qdrant Vector Database**: Standard SSD, 50GB. Fast disk recommended for AI embeddings.
- **NATS JetStream**: SSD, 20GB for message persistence.
- **Prometheus**: SSD, 30GB for metrics retention.

## 6. Networking Architecture
- **Private GKE Cluster**: Worker nodes have no public IP addresses.
- **Cloud NAT**: Allows outbound traffic from nodes to the internet (e.g., pulling images, calling external APIs).
- **Ingress Controller**: Google Cloud Load Balancer (GCLB) sitting in a public subnet, routing HTTPS traffic to the private pods via Network Endpoint Groups (NEGs).

## 7. DNS Architecture (Cloudflare)
The domain `cortexops.amshithnair.in` is managed in Cloudflare.
- **A Records**: Point to the external IP of the GCLB.
  - `cortexops.amshithnair.in` (Frontend -> Vercel, *see Monorepo guide*)
  - `api.cortexops.amshithnair.in` -> GCLB IP
  - `grafana.cortexops.amshithnair.in` -> GCLB IP
  - `temporal.cortexops.amshithnair.in` -> GCLB IP

## 8. TLS Architecture
- **Cloudflare Full (Strict) Mode**: Encrypts traffic from User -> Cloudflare -> GCP.
- **cert-manager**: Installed in GKE to automatically provision Let's Encrypt certificates using DNS-01 challenge (via Cloudflare API) or Cloudflare Origin CA certificates.
- **Ingress TLS**: The GCLB terminates TLS using the certificate managed by cert-manager and forwards traffic to the pods.

## 9. Cloudflare Configuration
1. Go to Cloudflare Dashboard -> SSL/TLS.
2. Set encryption mode to **Full (Strict)**.
3. In DNS, create the A records for `api`, `grafana`, `temporal` pointing to the Load Balancer IP. Make sure they are **Proxied (Orange Cloud)**.
4. Go to SSL/TLS -> Origin Server and generate an Origin Certificate.
5. Store this certificate as a Kubernetes Secret for your Ingress to use.

## 10. Artifact Registry Setup
1. Enable Artifact Registry API in GCP.
2. Create a Docker repository: `gcloud artifacts repositories create cortexops-repo --repository-format=docker --location=us-central1`
3. Configure Docker auth: `gcloud auth configure-docker us-central1-docker.pkg.dev`
4. Set up GCP Workload Identity or Service Account keys for GitHub Actions to push images.

## 11. GitHub Actions Deployment Workflow
- **Trigger**: Push to `main` branch.
- **Build**: Docker build for each Go microservice.
- **Push**: Push to `us-central1-docker.pkg.dev/PROJECT_ID/cortexops-repo/service:tag`.
- **Deploy**: Update Helm values with the new tag and run `helm upgrade --install cortexops deploy/helm/cortexops -f values.yaml`.

## 12. Rollback Strategy
- **Helm Rollback**: Use `helm history cortexops` and `helm rollback cortexops <revision>` for instant rollback of Kubernetes manifests and images.
- **Database Migrations**: Ensure schema changes are backwards compatible.

## 13. Disaster Recovery Strategy
- **Database**: Enable automated daily backups in Cloud SQL.
- **Kubernetes State**: Install **Velero** to backup cluster state (ConfigMaps, Secrets, Ingress) to a GCS bucket.
- **Qdrant**: Configure snapshotting and sync to a GCS bucket.

## 14. Step-by-step Deployment Instructions

### Phase 1: GCP Foundation
1. **Login**: `gcloud auth login`
2. **Project**: `gcloud config set project [YOUR_PROJECT_ID]`
3. **Enable APIs**: `gcloud services enable container.googleapis.com artifactregistry.googleapis.com sqladmin.googleapis.com`

### Phase 2: Cluster Creation
1. Create a VPC: `gcloud compute networks create cortexops-vpc --subnet-mode=auto`
2. Create Autopilot Cluster:
   ```bash
   gcloud container clusters create-auto cortexops-cluster \
     --region us-central1 \
     --network cortexops-vpc
   ```
3. Get credentials: `gcloud container clusters get-credentials cortexops-cluster --region us-central1`

### Phase 3: Dependencies Setup
1. **Install Nginx Ingress or use GCLB**:
   If using Nginx (recommended for complex routing):
   `helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx --namespace ingress-nginx --create-namespace`
2. **Install cert-manager**:
   `helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --set installCRDs=true`
3. **Install Prometheus/Grafana**:
   `helm upgrade --install monitoring prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace`

### Phase 4: Application Deployment
1. Build and push images to Artifact Registry.
2. Update `deploy/helm/cortexops/values.yaml` to point to your Artifact Registry images.
3. Install the application:
   ```bash
   helm upgrade --install cortexops deploy/helm/cortexops \
     --namespace cortexops --create-namespace \
     --set ingress.enabled=true \
     --set ingress.hosts[0].host=api.cortexops.amshithnair.in \
     --set ingress.hosts[0].paths[0].path=/ \
     --set ingress.tls[0].hosts[0]=api.cortexops.amshithnair.in
   ```
4. Verify deployment: `kubectl get pods -n cortexops`
