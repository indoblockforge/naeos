package pluginhost

import (
	"testing"
)

func TestEventBusGetter(t *testing.T) {
	m := NewManager("")
	if m.EventBus() == nil {
		t.Error("expected non-nil EventBus")
	}
}

func TestPluginEventBusNoOps(t *testing.T) {
	p := NewPluginEventBus(NewEventBus())
	p.OnPipelineStart("p1")
	p.OnArtifactGenerated("a", "/tmp/a")
}

func TestSimpleMetricsSnapshotEmpty(t *testing.T) {
	m := NewSimpleMetrics()
	counters, gauges, histograms := m.Snapshot()
	if len(counters) != 0 || len(gauges) != 0 || len(histograms) != 0 {
		t.Error("expected empty snapshot on fresh metrics")
	}
}
