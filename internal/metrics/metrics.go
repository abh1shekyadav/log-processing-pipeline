package metrics

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	generated  uint64
	processed  uint64
	aggregated uint64
	startTime  = time.Now()
)

func IncGenerated() {
	atomic.AddUint64(&generated, 1)
}

func IncProcessed() {
	atomic.AddUint64(&processed, 1)
}

func IncAggregated() {
	atomic.AddUint64(&aggregated, 1)
}

func Snapshot() {
	fmt.Printf("\n[METRICS] Stats since %v\n", startTime.Format("15:04:05"))
	fmt.Printf("  Logs Generated : %d\n", atomic.LoadUint64(&generated))
	fmt.Printf("  Logs Processed : %d\n", atomic.LoadUint64(&processed))
	fmt.Printf("  Logs Aggregated: %d\n", atomic.LoadUint64(&aggregated))
	fmt.Println("--------------------------------------------")
}

func Reset() {
	atomic.StoreUint64(&generated, 0)
	atomic.StoreUint64(&processed, 0)
	atomic.StoreUint64(&aggregated, 0)
	startTime = time.Now()
}
