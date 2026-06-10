export default function ReferencePage() {
  return (
    <>
      <h1 className="text-4xl font-bold text-white mb-4">Reference</h1>
      <p className="text-xl text-zinc-400 mb-12 font-light">
        API definitions, commands, and state definitions.
      </p>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Diagnostics API</h2>
        <div className="bg-zinc-900/50 border border-zinc-800 rounded-xl overflow-hidden">
          <table className="w-full text-left text-sm text-zinc-300">
            <thead className="bg-zinc-800 text-zinc-100">
              <tr>
                <th className="px-6 py-3 font-medium">Endpoint</th>
                <th className="px-6 py-3 font-medium">Description</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800">
              <tr>
                <td className="px-6 py-4 font-mono text-cortex-400">GET /debug/healthz</td>
                <td className="px-6 py-4">Returns 200 if NATS and PG are reachable.</td>
              </tr>
              <tr>
                <td className="px-6 py-4 font-mono text-cortex-400">GET /debug/incidents/active</td>
                <td className="px-6 py-4">Lists currently tracking incidents in the CE.</td>
              </tr>
              <tr>
                <td className="px-6 py-4 font-mono text-cortex-400">GET /debug/graph/blast-radius</td>
                <td className="px-6 py-4">Returns full dependency tree for a given NodeID.</td>
              </tr>
              <tr>
                <td className="px-6 py-4 font-mono text-cortex-400">GET /debug/temporal/status</td>
                <td className="px-6 py-4">Checks workflow queue latency.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Make Commands</h2>
        <ul className="space-y-4">
          <li className="bg-black border border-zinc-800 p-4 rounded-lg">
            <code className="text-cortex-400 font-mono block mb-2">make dev-up</code>
            <p className="text-sm text-zinc-400">Spins up the stateful dependencies (Postgres, NATS, Temporal) locally via Docker Compose.</p>
          </li>
          <li className="bg-black border border-zinc-800 p-4 rounded-lg">
            <code className="text-cortex-400 font-mono block mb-2">make verify-runtime</code>
            <p className="text-sm text-zinc-400">Validates that all local systems are online and reachable.</p>
          </li>
          <li className="bg-black border border-zinc-800 p-4 rounded-lg">
            <code className="text-cortex-400 font-mono block mb-2">make bootstrap</code>
            <p className="text-sm text-zinc-400">Deploys the demo microservices into the `cortexops-demo` namespace.</p>
          </li>
          <li className="bg-black border border-zinc-800 p-4 rounded-lg">
            <code className="text-cortex-400 font-mono block mb-2">make demo-failure SCENARIO=[scenario]</code>
            <p className="text-sm text-zinc-400">Injects a fault to trigger the correlation engine. Supported: `rollout-fail`, `crashloop`, `scale-pressure`.</p>
          </li>
        </ul>
      </section>

      <section className="mb-16">
        <h2 className="text-2xl font-semibold text-white mb-6 border-b border-zinc-800 pb-2">Workflow States</h2>
        <ul className="list-disc pl-5 text-sm text-zinc-300 space-y-2">
          <li><strong>PROPOSED:</strong> Incident correlated, RCA generated, workflow initiated.</li>
          <li><strong>POLICY_EVALUATING:</strong> OPA rego rule evaluation.</li>
          <li><strong>APPROVAL_PENDING:</strong> Blocked waiting for Human-in-the-loop (Slack).</li>
          <li><strong>DRY_RUNNING:</strong> Validating K8s API patches without mutating state.</li>
          <li><strong>EXECUTING:</strong> Applying mutations.</li>
          <li><strong>VERIFYING:</strong> Monitoring telemetry post-execution for stability.</li>
          <li><strong>SUCCESS:</strong> Remediation resolved the incident.</li>
          <li><strong>ROLLING_BACK:</strong> Verifying failed, executing compensation transactions.</li>
          <li><strong>REJECTED:</strong> Failed OPA evaluation or human rejected.</li>
        </ul>
      </section>
    </>
  );
}
