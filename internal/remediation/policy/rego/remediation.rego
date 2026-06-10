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
}

action_allowed {
    allowed_actions[input.action.type]
}

namespace_protected {
    protected_namespaces[input.action.namespace]
}

# Approval requirement rule
requires_approval {
    input.risk_score > 0.5
}
