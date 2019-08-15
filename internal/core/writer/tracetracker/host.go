package tracetracker

import (
	"github.com/signalfx/golib/trace"
	"github.com/signalfx/signalfx-agent/internal/core/common/dpmeta"
	"github.com/signalfx/signalfx-agent/internal/core/services"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
)

var dimsToSyncSource = []string{
	"container_id",
	"kubernetes_pod_uid",
}

type SourceTracker struct {
	dimensionCh chan<- *types.Dimension
}

func NewSourceTracker(dimensionCh chan<- *types.Dimension) *SourceTracker {
	return &SourceTracker{
		dimensionCh: dimensionCh,
	}
}

func (st *SourceTracker) AddSourceTagsToSpan(span *trace.Span) {
	endpoint, ok := span.Meta[dpmeta.EndpointMeta].(services.Endpoint)
	if !ok || endpoint == nil {
		return
	}

	dims := endpoint.Dimensions()
	for _, dim := range dimsToSyncSource {
		if val := dims[dim]; val != "" {
			tags := span.Tags
			if _, ok := tags[dim]; ok {
				// Don't overwrite existing span tags
				continue
			}
			tags[dim] = val
		}
	}
}
