export default function DeploymentGuidePage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Deploy CortexOps</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        Deploy CortexOps on Kubernetes in minutes.
      </p>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Supported Platforms</h2>
        <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2 mb-6">
          <li><strong>Kubernetes</strong> (Vanilla, 1.28+)</li>
          <li><strong>GKE</strong> (Google Kubernetes Engine)</li>
          <li><strong>EKS</strong> (Elastic Kubernetes Service)</li>
          <li><strong>AKS</strong> (Azure Kubernetes Service)</li>
          <li><strong>Kind</strong> (Kubernetes in Docker, for local testing)</li>
        </ul>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Deployment Methods</h2>
        <p className="text-zinc-300 mb-4">
          CortexOps can be deployed using either Helm Charts (recommended for production) or pure Kubernetes Manifests (for evaluation).
        </p>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Requirements</h2>
        <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2 mb-6">
          <li>kubectl installed and configured</li>
          <li>Helm v3.x installed</li>
          <li>At least 4 CPU cores and 8GB RAM available in the cluster</li>
          <li>Storage class supporting persistent volume claims (for PostgreSQL and NATS)</li>
        </ul>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Installation via Helm</h2>
        <div className="bg-black border border-zinc-800 rounded-lg p-4 overflow-x-auto text-sm text-zinc-300 mb-6 font-mono">
          # Add the CortexOps Helm repository
          helm repo add cortexops https://charts.cortexops.amshithnair.in
          helm repo update
          
          # Install the control plane
          helm install cortexops cortexops/control-plane \
            --namespace cortexops-system \
            --create-namespace
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Verification</h2>
        <p className="text-zinc-300 mb-4">
          Ensure all components are running correctly before injecting telemetry:
        </p>
        <div className="bg-black border border-zinc-800 rounded-lg p-4 overflow-x-auto text-sm text-zinc-300 mb-6 font-mono">
          kubectl get pods -n cortexops-system
          
          # You should see:
          # - cortexops-collector
          # - cortexops-correlator
          # - nats-cluster
          # - temporal-server
          # - postgresql
          # - qdrant
        </div>
      </section>
      
      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Troubleshooting</h2>
        <p className="text-zinc-300 mb-4">
          If pods are failing to start, check the operator logs or ensure your cluster has adequate resources allocated to the `cortexops-system` namespace.
        </p>
        <div className="bg-zinc-900/50 border-l-4 border-cortex-500 p-6 rounded-r-xl mb-6 text-zinc-300 italic">
          For detailed diagnostic commands, refer to the Operations Runbooks.
        </div>
      </section>
    </>
  );
}
