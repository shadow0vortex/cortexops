package cortexops.remediation

default allowed = false

# Allow specific action types
allowed {
    valid_action_type[input.action.type]
    count(deny) == 0
}

valid_action_type := {
    "POD_RESTART",
    "DEPLOYMENT_ROLLOUT_RESTART",
    "HORIZONTAL_SCALE"
}

# Deny mutation in kube-system
deny[reason] {
    input.action.target_namespace == "kube-system"
    reason := "kube-system is immutable to automation"
}

# Deny if invalid action
deny[reason] {
    not valid_action_type[input.action.type]
    reason := sprintf("Action type %v is strictly forbidden", [input.action.type])
}
