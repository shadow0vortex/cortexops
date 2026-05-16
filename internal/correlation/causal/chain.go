package causal

import (
	"context"
	"sort"

	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	eventsv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/pkg/core"
)

// ChainBuilder constructs deterministic propagation chains from a bucket of correlated events.
type ChainBuilder struct {
	topology core.TopologyProvider
}

func NewChainBuilder(topology core.TopologyProvider) *ChainBuilder {
	return &ChainBuilder{topology: topology}
}

// Build orders events chronologically and assigns causal links based on topology directed edges.
func (c *ChainBuilder) Build(ctx context.Context, evidence []*eventsv1.TelemetryEnvelope) (*correlationv1.CausalChain, error) {
	if len(evidence) == 0 {
		return nil, nil
	}

	// 1. Sort events chronologically to establish temporal precedence
	sort.Slice(evidence, func(i, j int) bool {
		return evidence[i].Timestamp.Seconds < evidence[j].Timestamp.Seconds
	})

	rootCauseEvent := evidence[0]
	var links []*correlationv1.CausalLink

	// 2. Iterate and attempt to link each event to a preceding event using Topology
	for i := 1; i < len(evidence); i++ {
		current := evidence[i]
		
		// Look backwards chronologically to find the closest topological parent
		var linked bool
		for j := i - 1; j >= 0; j-- {
			preceding := evidence[j]
			
			// We only build topology links for K8s events right now
			if current.GetK8SEvent() != nil && preceding.GetK8SEvent() != nil {
				curK := current.GetK8SEvent()
				preK := preceding.GetK8SEvent()
				
				curID := curK.ResourceKind + "/" + curK.Namespace + "/" + curK.ResourceName
				preID := preK.ResourceKind + "/" + preK.Namespace + "/" + preK.ResourceName
				
				// Does Preceding Event's Node have an outward edge to Current Event's Node?
				deps, _ := c.topology.GetDependencies(ctx, preID)
				for _, dep := range deps {
					if dep == curID {
						links = append(links, &correlationv1.CausalLink{
							CauseEventId:  preceding.EventId,
							EffectEventId: current.EventId,
							Reasoning:     "Topology propagation (" + preID + " -> " + curID + ")",
						})
						linked = true
						break
					}
				}
			}
			if linked {
				break
			}
		}
		
		// If no topology link was found, default to a temporal sequence link
		if !linked {
			links = append(links, &correlationv1.CausalLink{
				CauseEventId:  evidence[i-1].EventId,
				EffectEventId: current.EventId,
				Reasoning:     "Temporal precedence",
			})
		}
	}

	return &correlationv1.CausalChain{
		RootCauseEventId: rootCauseEvent.EventId,
		Links:            links,
	}, nil
}
