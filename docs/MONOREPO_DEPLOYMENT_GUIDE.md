# CortexOps Monorepo Deployment Guide

This guide outlines the CI/CD and deployment strategies for the CortexOps monorepo. It focuses on isolating frontend and backend workflows and deploying them to their optimized target environments.

## Repository Structure Context
- `frontend/`: Next.js React application.
- `cmd/`, `internal/`, `pkg/`: Go backend microservices (Collector, Correlator, Topology, RCA, Remediation).
- `deploy/helm/`: Kubernetes manifests.
- `docs/`: Documentation.

## 1. Vercel Deployment for Frontend Only
Vercel is the optimal hosting platform for Next.js. We will decouple the frontend from the GKE backend.

1. Connect Vercel to the GitHub repository.
2. During the project import, under **Framework Preset**, select `Next.js`.

## 2. Correct Root Directory Configuration
Crucially, tell Vercel where the Next.js app lives within the monorepo:
- **Root Directory**: `frontend`
- Do not include the trailing slash. Vercel will run `npm install` and `npm run build` from inside this directory.

## 3. Next.js Environment Variables
In the Vercel project settings, configure the following Environment Variables:
- `NEXT_PUBLIC_API_URL`: `https://api.cortexops.amshithnair.in` (Points to the GKE backend).
- `NEXT_PUBLIC_WS_URL`: `wss://api.cortexops.amshithnair.in/ws` (If using WebSockets).

## 4. Cloudflare DNS Setup
Manage your domain (`cortexops.amshithnair.in`) in Cloudflare.

| Hostname | Record Type | Target | Proxy Status |
| :--- | :--- | :--- | :--- |
| `cortexops.amshithnair.in` | CNAME | `cname.vercel-dns.com` | Proxied (Orange Cloud) |
| `api.cortexops.amshithnair.in` | A | `<GKE_INGRESS_IP>` | Proxied (Orange Cloud) |
| `grafana.cortexops.amshithnair.in` | A | `<GKE_INGRESS_IP>` | Proxied (Orange Cloud) |
| `temporal.cortexops.amshithnair.in` | A | `<GKE_INGRESS_IP>` | Proxied (Orange Cloud) |

*(Vercel may require you to unproxy the CNAME temporarily during initial SSL verification).*

## 5. GKE Deployment Architecture (Backend)
The backend runs on GKE Autopilot. 
- An **Ingress Controller** (e.g., NGINX Ingress or GKE Ingress) listens on port 80/443 via a LoadBalancer service.
- The Ingress resource uses host-based routing to direct traffic to the correct internal service.

## 6. Backend Subdomain Architecture
The Kubernetes Ingress rules map subdomains to services:

- `api.cortexops.amshithnair.in/*` -> routes to `cortexops-api-gateway` or directly to `cortexops-collector`/`cortexops-correlator` depending on path routing configuration.
- `grafana.cortexops.amshithnair.in/*` -> routes to `kube-prometheus-stack-grafana` service on port 80.
- `temporal.cortexops.amshithnair.in/*` -> routes to `temporal-web` service on port 8080.

## 7. CI/CD Strategy for Monorepo
To save GitHub Actions minutes and prevent unnecessary deployments, use **path filtering**. 
- A change in `frontend/` should only trigger the Vercel deployment.
- A change in `cmd/` or `internal/` should trigger the Go build and GKE deployment.

## 8. GitHub Actions Workflow Separation

Create separate workflow files in `.github/workflows/`:

### `.github/workflows/frontend.yaml`
*(Optional: Vercel's native GitHub integration handles deployment automatically, but you can run tests here).*
```yaml
name: Frontend CI
on:
  push:
    branches: [ main ]
    paths:
      - 'frontend/**'
jobs:
  test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./frontend
    steps:
      - uses: actions/checkout@v4
      - run: npm ci
      - run: npm run lint
      - run: npm run build
```

### `.github/workflows/backend.yaml`
```yaml
name: Backend Deploy to GKE
on:
  push:
    branches: [ main ]
    paths:
      - 'cmd/**'
      - 'internal/**'
      - 'pkg/**'
      - 'deploy/**'
      - 'build/docker/**'
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Authenticate to Google Cloud
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_CREDENTIALS }}
      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v2
      - name: Configure Docker
        run: gcloud auth configure-docker us-central1-docker.pkg.dev
      - name: Build & Push Images
        run: |
          # Build Collector
          docker build -f build/docker/Dockerfile.base -t us-central1-docker.pkg.dev/PROJECT/cortexops-repo/collector:latest .
          docker push us-central1-docker.pkg.dev/PROJECT/cortexops-repo/collector:latest
          # Repeat for other services...
      - name: Get GKE Credentials
        uses: google-github-actions/get-gke-credentials@v2
        with:
          cluster_name: cortexops-cluster
          location: us-central1
      - name: Deploy with Helm
        run: |
          helm upgrade --install cortexops deploy/helm/cortexops -n cortexops --create-namespace
```

## Step-by-Step Monorepo Deployment Instructions for Beginners

1. **Frontend (Vercel)**
   - Create an account at [vercel.com](https://vercel.com).
   - Click "Add New" -> "Project".
   - Import your GitHub repository.
   - Expand "Build and Output Settings".
   - Crucially: Change **Root Directory** to `frontend`.
   - Add Environment Variable `NEXT_PUBLIC_API_URL` = `https://api.cortexops.amshithnair.in`.
   - Click Deploy.
   - Once deployed, add `cortexops.amshithnair.in` to the Vercel Custom Domains setting.

2. **Backend (GCP & GitHub Actions)**
   - Follow the `GKE_DEPLOYMENT_BLUEPRINT.md` to set up your GKE cluster.
   - In GitHub -> Settings -> Secrets and variables -> Actions, add your `GCP_CREDENTIALS` as a secret.
   - Push your code. The GitHub Actions workflow (defined in section 8) will trigger because it detects changes in the backend folders.
   - The workflow will build your Go binaries, package them in Docker containers, push them to Artifact Registry, and run `helm upgrade` on your GKE cluster.

3. **DNS Hookup (Cloudflare)**
   - Go to Cloudflare -> DNS.
   - Create a CNAME record for `cortexops` pointing to `cname.vercel-dns.com`.
   - Get your Ingress External IP: `kubectl get ingress -n cortexops`.
   - Create A records for `api`, `grafana`, and `temporal` pointing to that IP address.
   - Wait for DNS propagation and test your fully deployed monorepo application!
