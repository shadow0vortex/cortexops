package remediation.policy

default allow = false

# Allowed Action Types
allowed_actions = {
    "POD_RESTART",
    "DEPLOYMENT_ROLLOUT_RESTART",
    "HORIZONTAL_SCALE"
}

# Protected Namespaces
protected_namespaces = {
    "kube-system",
    "kube-public",
    "kube-node-lease",
    "cortexops"
}

# Main allow rule
allow {
    action_allowed
    not namespace_protected
    not maintenance_active
}

action_allowed {
    allowed_actions[input.action.type]
}

namespace_protected {
    protected_namespaces[input.action.namespace]
}

# Maintenance Window Guard (TD-031)
# Blocks all automated remediation when a maintenance window is active.
# The caller must inject input.maintenance_window = true when the system
# is within a scheduled maintenance period.
maintenance_active {
    input.maintenance_window == true
}

# Denial reasons for audit trail
deny[reason] {
    not action_allowed
    reason := sprintf("Action type '%v' is not in the allowed set", [input.action.type])
}

deny[reason] {
    namespace_protected
    reason := sprintf("Namespace '%v' is protected from automated remediation", [input.action.namespace])
}

deny[reason] {
    maintenance_active
    reason := "Remediation blocked: maintenance window is active"
}

# Approval requirement rule
requires_approval {
    input.risk_score > 0.5
}
