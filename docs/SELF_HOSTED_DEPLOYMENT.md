# Self-Hosted Production Deployment

This runbook deploys the five Go services and their stateful dependencies to a Kubernetes cluster on your server. The Next.js frontend remains on Vercel. Use a dedicated API hostname, such as `api.example.com`, for the cluster ingress; point the frontend hostname to Vercel through the Vercel domain settings.

## Prerequisites

- A supported Kubernetes cluster on the server, `kubectl`, and Helm 3.14 or newer.
- A default StorageClass or an explicitly selected class with persistent-volume support.
- A Docker Hub account able to push to `shadow0vortex/cortexops-{service}` repositories.
- DNS control for the domain and public server access on TCP ports 80 and 443.
- An email address for Let's Encrypt certificate registration.

The included PostgreSQL, Qdrant, and NATS charts create PVCs. Temporal stores workflow history in PostgreSQL and does not need a separate database volume. Select PVC sizes and StorageClass for the server before installation; the defaults are starting values, not capacity guarantees.

## Publish Backend Images

Add these GitHub Actions repository secrets:

- `DOCKERHUB_USERNAME`: Docker Hub username.
- `DOCKERHUB_TOKEN`: Docker Hub access token with write permission.

Pull requests build images without publishing. Pushes to `main` publish commit and `latest` tags; version tags publish versioned images. All five Go service images use a shared tag:

```text
docker.io/shadow0vortex/cortexops-collector:<tag>
docker.io/shadow0vortex/cortexops-correlator:<tag>
docker.io/shadow0vortex/cortexops-topology:<tag>
docker.io/shadow0vortex/cortexops-rca:<tag>
docker.io/shadow0vortex/cortexops-remediation:<tag>
```

Use the same immutable commit or release tag for `image.tag` during deployment. Do not deploy `latest` for a production rollout.

## Install Ingress and TLS

Install the controller and cert-manager once per cluster. The chart ingress uses class `nginx` and the cert-manager ClusterIssuer name `letsencrypt-prod` by default.

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx --create-namespace
helm upgrade --install cert-manager jetstack/cert-manager \
  --namespace cert-manager --create-namespace --set crds.enabled=true
```

Create a ClusterIssuer after cert-manager is ready. Replace the email with an address monitored by your team:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    email: ops@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod-account
    solvers:
      - http01:
          ingress:
            class: nginx
```

Apply it with `kubectl apply -f cluster-issuer.yaml`. Ensure the API hostname resolves to the ingress controller's external IP before expecting certificate issuance.

## Create Runtime Secrets

Create the namespace, then provide secrets before installing the chart. Use a secret manager or external-secrets integration for production; do not commit secrets or put them in Helm values. The following Bash example is for a manually managed cluster secret setup:

```bash
kubectl create namespace cortexops

read -rsp 'Database password: ' POSTGRES_PASSWORD; echo
read -rsp 'NATS password: ' NATS_PASSWORD; echo
read -rsp 'Diagnostics API token: ' DIAG_API_TOKEN; echo
NATS_USER=cortexops

kubectl -n cortexops create secret generic cortexops-postgres-auth \
  --from-literal=password="$POSTGRES_PASSWORD"
kubectl -n cortexops create secret generic cortexops-postgres \
  --from-literal=url="postgres://cortex:${POSTGRES_PASSWORD}@cortexops-postgresql:5432/cortexops?sslmode=disable"
kubectl -n cortexops create secret generic cortexops-nats-auth \
  --from-literal=username="$NATS_USER" \
  --from-literal=password="$NATS_PASSWORD" \
  --from-literal=url="nats://${NATS_USER}:${NATS_PASSWORD}@cortexops-nats:4222"
kubectl -n cortexops create secret generic cortexops-diag-token \
  --from-literal=token="$DIAG_API_TOKEN"
unset POSTGRES_PASSWORD NATS_PASSWORD DIAG_API_TOKEN
```

The chart expects `cortexops-postgres-auth` key `password`; `cortexops-postgres` key `url`; `cortexops-nats-auth` keys `username`, `password`, and `url`; and `cortexops-diag-token` key `token`. Configure secret rotation with the chosen secret-management system and roll deployments after rotating values.

## Kubernetes Access

`clusterWideAccess` defaults to `false`, limiting Kubernetes observation to the release namespace. Set it to `true` when the collector must monitor resources across the cluster; this grants cluster-wide `get`, `list`, and `watch` permissions to the collector and remediation service accounts.

`allowRemediationWrites` defaults to `false`. Keep it disabled for monitoring-only deployments. Enabling it grants the remediation service account permission to delete pods and patch deployments at the selected access scope (`clusterWideAccess` determines whether that scope is namespace-wide or cluster-wide).

The [homelab values profile](../deploy/helm/cortexops/values-homelab.yaml) enables cluster-wide read-only observation, uses one replica per CortexOps service, and disables autoscaling and pod disruption budgets. The [Argo CD Application](../deploy/argocd/application.yaml) uses this profile and requires a manual sync; it does not automatically prune resources. It has no cascading deletion finalizer, so deleting the Application leaves its resources in place for deliberate cleanup.

## Install CortexOps

The chart includes the pinned NATS, PostgreSQL, Qdrant, and Temporal dependencies. The root chart's PostgreSQL image is pinned by digest for both amd64 and arm64. Register the upstream repositories, build the committed lock file, then install using the same backend image tag across all services:

```bash
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add qdrant https://qdrant.github.io/qdrant-helm
helm repo add temporal https://go.temporal.io/helm-charts
helm repo update
helm dependency build deploy/helm/cortexops

helm upgrade --install cortexops deploy/helm/cortexops \
  --namespace cortexops --create-namespace \
  --set image.tag=<immutable-image-tag> \
  --set ingress.enabled=true \
  --set ingress.host=api.example.com
```

The chart exposes only the topology `/v1` API route through ingress. Health, metrics, Grafana, Temporal UI, NATS monitoring, databases, and Qdrant remain private to the cluster. The `POSTGRES_URL` and `DIAG_API_TOKEN` are required from Kubernetes Secrets by the topology service.

## Configure DNS and Vercel

1. Point the frontend hostname (for example, `example.com`) to the Vercel project using the DNS values shown in Vercel's domain settings.
2. Point the API hostname (for example, `api.example.com`) to the ingress controller's external IP. If the server IP changes, use a stable address or update DNS accordingly.
3. Set the Vercel project custom domain and verify its certificate in Vercel.
4. If the frontend needs to call CortexOps, configure its public API URL in Vercel's Production environment and redeploy. The current frontend source does not contain a clear API base-URL integration, so verify the intended frontend/API contract before depending on it.
5. Keep the diagnostics bearer token in the appropriate server-side secret store. Do not expose it as a `NEXT_PUBLIC_*` variable, which would ship it to browsers.

The actual domain, DNS records, ingress IP, ACME email, and any Cloudflare proxy settings are environment-specific and are not changed by this repository.

## Verify, Upgrade, and Recover

```bash
kubectl -n cortexops get pods,pvc,svc,ingress
kubectl -n cortexops get certificate
helm -n cortexops status cortexops
helm -n cortexops history cortexops
```

Check that PVCs are bound, all five Go Deployments are ready, the certificate is `Ready`, and `https://api.example.com/v1/topology/nodes` returns `401` without a bearer token. An authenticated request should return `200`.

For upgrades, use a newly published immutable backend image tag and run `helm upgrade --install` again with the same ingress values. Roll back a release with `helm rollback cortexops <revision> -n cortexops`; validate PostgreSQL and Qdrant backups before relying on PVC retention as disaster recovery. PVCs protect against pod replacement, not disk/server loss. Test restoring database and Qdrant snapshots on a separate volume before production use.