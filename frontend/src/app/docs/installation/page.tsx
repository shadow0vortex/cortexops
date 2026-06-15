import { Terminal, HardDrive, Server, Cloud } from "lucide-react";

export default function InstallationPage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Installation</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        Procedures to install, configure, and run CortexOps across various environments.
      </p>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Prerequisites</h2>
        <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl mb-6 text-sm text-zinc-300">
          <p className="mb-3">Before installing CortexOps, ensure your environment meets the following requirements:</p>
          <ul className="list-disc pl-5 space-y-2">
            <li><strong>Docker & Docker Compose:</strong> Required for local development and testing.</li>
            <li><strong>Kubernetes Cluster:</strong> (v1.26+) For production deployments (e.g., Minikube, kind, GKE, EKS).</li>
            <li><strong>kubectl:</strong> Configured to communicate with your cluster.</li>
            <li><strong>Helm:</strong> (v3.0+) For deploying the CortexOps charts.</li>
            <li><strong>GNU Make:</strong> To run deployment commands via the provided Makefile.</li>
          </ul>
        </div>
      </section>

      <section className="mb-16">
        <div className="flex items-center gap-3 mb-6 border-b border-zinc-800 pb-2">
          <HardDrive className="w-6 h-6 text-cortex-400" />
          <h2 className="text-2xl font-semibold text-white">Local Development (Docker Compose)</h2>
        </div>
        <p className="text-zinc-300 mb-6">
          The fastest way to run CortexOps on your local machine is using Docker Compose. This spins up the stateful dependencies (Postgres, NATS, Temporal) locally.
        </p>

        <div className="space-y-6">
          <div>
            <h3 className="text-lg font-medium text-white mb-2">1. Clone the Repository</h3>
            <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300">
              <div className="flex items-center gap-3">
                <Terminal className="w-4 h-4 text-zinc-500" />
                <span>git clone https://github.com/cortexops/cortexops.git<br/>cd cortexops</span>
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-medium text-white mb-2">2. Start Infrastructure</h3>
            <p className="text-sm text-zinc-400 mb-3">Initialize all backing services in detached mode.</p>
            <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300">
              <div className="flex items-center gap-3">
                <Terminal className="w-4 h-4 text-zinc-500" />
                <span>make dev-up</span>
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-medium text-white mb-2">3. Verify the Runtime</h3>
            <p className="text-sm text-zinc-400 mb-3">Ensure that NATS, Temporal, and PostgreSQL are healthy.</p>
            <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300">
              <div className="flex items-center gap-3">
                <Terminal className="w-4 h-4 text-zinc-500" />
                <span>make verify-runtime</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="mb-16">
        <div className="flex items-center gap-3 mb-6 border-b border-zinc-800 pb-2">
          <Server className="w-6 h-6 text-cortex-400" />
          <h2 className="text-2xl font-semibold text-white">Kubernetes Installation (Helm)</h2>
        </div>
        <p className="text-zinc-300 mb-6">
          Deploy CortexOps into an existing Kubernetes cluster for production workloads or staging environments.
        </p>

        <div className="space-y-6">
          <div>
            <h3 className="text-lg font-medium text-white mb-2">1. Add the Helm Repository</h3>
            <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300">
              <div className="flex items-center gap-3">
                <Terminal className="w-4 h-4 text-zinc-500" />
                <span>helm repo add cortexops https://charts.cortexops.io<br/>helm repo update</span>
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-medium text-white mb-2">2. Create Namespace</h3>
            <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300">
              <div className="flex items-center gap-3">
                <Terminal className="w-4 h-4 text-zinc-500" />
                <span>kubectl create namespace cortexops</span>
              </div>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-medium text-white mb-2">3. Deploy the Chart</h3>
            <p className="text-sm text-zinc-400 mb-3">Deploy the control plane using the default configuration values.</p>
            <div className="relative bg-black border border-zinc-800 rounded-lg p-4 flex items-center justify-between font-mono text-sm text-zinc-300">
              <div className="flex items-center gap-3">
                <Terminal className="w-4 h-4 text-zinc-500" />
                <span>helm install cortexops cortexops/cortexops -n cortexops</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="mb-16">
        <div className="flex items-center gap-3 mb-6 border-b border-zinc-800 pb-2">
          <Cloud className="w-6 h-6 text-cortex-400" />
          <h2 className="text-2xl font-semibold text-white">GitOps & Managed Environments</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl">
            <h3 className="text-lg font-medium text-white mb-2">ArgoCD / Flux</h3>
            <p className="text-sm text-zinc-400">
              We highly recommend deploying CortexOps via GitOps controllers. You can find pre-configured <code>Application</code> manifests in the <code>/deploy/gitops</code> directory of the repository.
            </p>
          </div>
          <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl">
            <h3 className="text-lg font-medium text-white mb-2">Cloud Providers (GKE / EKS)</h3>
            <p className="text-sm text-zinc-400">
              For managed Kubernetes platforms, ensure that your node pools have sufficient capacity (minimum 4vCPU and 16GB RAM total) to host the stateful layer (Temporal + Postgres).
            </p>
          </div>
        </div>
      </section>
    </>
  );
}
