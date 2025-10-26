package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type PipelineMetrics struct {
	Name       string
	generated  uint64
	processed  uint64
	batches    uint64
	aggregated uint64
	start      time.Time
}

func NewPipelineMetrics(name string) *PipelineMetrics {
	return &PipelineMetrics{
		Name:  name,
		start: time.Now(),
	}
}

func (m *PipelineMetrics) IncGenerated()  { atomic.AddUint64(&m.generated, 1) }
func (m *PipelineMetrics) IncProcessed()  { atomic.AddUint64(&m.processed, 1) }
func (m *PipelineMetrics) IncBatches()    { atomic.AddUint64(&m.batches, 1) }
func (m *PipelineMetrics) IncAggregated() { atomic.AddUint64(&m.aggregated, 1) }

func (m *PipelineMetrics) Snapshot() map[string]interface{} {
	return map[string]interface{}{
		"name":       m.Name,
		"uptime_sec": time.Since(m.start).Seconds(),
		"generated":  atomic.LoadUint64(&m.generated),
		"processed":  atomic.LoadUint64(&m.processed),
		"batches":    atomic.LoadUint64(&m.batches),
		"aggregated": atomic.LoadUint64(&m.aggregated),
	}
}

func StartHTTPReporter(addr string, metricsList []*PipelineMetrics) {
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		out := map[string]interface{}{}
		for _, m := range metricsList {
			out[m.Name] = m.Snapshot()
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	})
	go func() {
		_ = http.ListenAndServe(addr, nil)
	}()
}

func (m *PipelineMetrics) Print() {
	s := m.Snapshot()
	b, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(b))
}
