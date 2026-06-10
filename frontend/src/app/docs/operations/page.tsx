export default function OperationsPage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Operations</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        First-class operational procedures for managing CortexOps in production.
      </p>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Runtime Operations</h2>
        <div className="space-y-4 text-sm text-zinc-300">
          <p><strong>Telemetry Flow:</strong> Monitor ingestion rates via the <code>cortexops_events_ingested_total</code> metric in Prometheus. Backpressure alerts indicate downstream NATS constraints.</p>
          <p><strong>Dashboard Metrics:</strong> The Grafana dashboard surfaces Correlation Latency, Workflow Success Rates, and LLM Token Usage.</p>
          <p><strong>Workflow States:</strong> Monitored via the Temporal UI. Look for workflows stuck in <code>APPROVAL_PENDING</code> or looping in <code>DRY_RUNNING</code>.</p>
          <p><strong>Runtime Visibility:</strong> Use the Diagnostics API (`/debug/healthz`) to verify sub-component availability.</p>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Incident Response</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl">
            <h3 className="text-lg font-medium text-white mb-2">High Ingestion Lag</h3>
            <p className="text-xs text-zinc-400">Occurs during massive cascading outages. CortexOps will auto-scale the Correlation Engine. If NATS hits memory limits, it drops older telemetry in favor of fresh signals.</p>
          </div>
          <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl">
            <h3 className="text-lg font-medium text-white mb-2">Failed RCA Generation</h3>
            <p className="text-xs text-zinc-400">If the LLM provider times out, CortexOps falls back to deterministic heuristic outputs. No workflows are blocked, but advisories become generic.</p>
          </div>
          <div className="bg-zinc-900/50 border border-zinc-800 p-6 rounded-xl">
            <h3 className="text-lg font-medium text-white mb-2">Stalled Workflows</h3>
            <p className="text-xs text-zinc-400">Workflows stuck in remediation usually indicate missing K8s RBAC permissions for the worker. Check Temporal activity logs for <code>Unauthorized</code> errors.</p>
          </div>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Remediation Operations</h2>
        <div className="space-y-4 text-sm text-zinc-300">
          <p><strong>Dry-Run Execution:</strong> All mutating actions are executed in dry-run mode against the API server first to validate syntax and permissions.</p>
          <p><strong>Rollback Handling:</strong> If post-remediation metrics do not stabilize within the configured window, the workflow executes a compensation transaction to revert the change.</p>
          <p><strong>Human Approval Workflows:</strong> High-risk operations (e.g., database restarts) pause execution and send a Slack interactive message for approval. They timeout after 15 minutes.</p>
          <p><strong>OPA Denials:</strong> Rejected operations log a policy violation event in the Audit DB and terminate immediately.</p>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Temporal Recovery</h2>
        <p className="text-sm text-zinc-400 mb-6">Temporal is the brain of our remediation orchestration. If it fails, all autonomous actions stall.</p>
        <div className="bg-red-950/20 border border-red-900/30 p-6 rounded-xl text-sm">
          <h3 className="text-lg font-medium text-red-400 mb-3">Playbook: Temporal Outage</h3>
          <ul className="list-disc pl-5 text-zinc-300 space-y-2">
            <li><strong>Database Failures:</strong> If Temporal loses PostgreSQL connection, check <code>docker compose logs postgres</code>. Ensure migrations have completed.</li>
            <li><strong>Worker Crashes:</strong> If a remediation worker dies, reboot the pod. Temporal will assign the open workflow task to the new worker seamlessly.</li>
            <li><strong>Workflow Termination:</strong> If a workflow is acting destructively, use <code>tctl workflow terminate -w [workflow_id]</code> to forcefully halt it.</li>
          </ul>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Chaos & Backup</h2>
        <div className="space-y-4 text-sm text-zinc-300">
          <p><strong>Chaos Testing:</strong> Run <code>make demo-failure</code> to simulate NATS outages, pod deletions, and rollback scenarios. Validate that exactly-once semantics hold.</p>
          <p><strong>Backup & Restore:</strong> Take daily pg_dump backups of PostgreSQL for Temporal state and Audit history. Qdrant volumes should be snapshot weekly.</p>
        </div>
      </section>
    </>
  );
}
