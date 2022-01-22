package balance

import (
	"fmt"
	"sync"

	"github.com/smallnest/weighted"
)

// WeightedRR ...
type WeightedRR struct {
	mux sync.Mutex
	rrw *weighted.RRW
}

func (w *WeightedRR) OnChange(c *Config) {
	if c == nil || len(c.Items) == 0 {
		fmt.Println("not change at all")
		return
	}
	w.rrw.RemoveAll()
	for _, item := range c.Items {
		w.Add(item.Name, item.Weight)
	}
}

func (w *WeightedRR) Reset() {
	panic("implement me")
}

func (w *WeightedRR) RemoveAll() {
	panic("implement me")
}

func (w *WeightedRR) Add(item interface{}, weight int) {
	w.rrw.Add(item, weight)
}

func (w *WeightedRR) All() map[interface{}]int {
	return w.rrw.All()
}

// Next ...
func (w *WeightedRR) Next() (item interface{}) {
	w.mux.Lock()
	defer w.mux.Unlock()
	return w.rrw.Next()

}

func NewWeightedRR(c *Config) *WeightedRR {
	rrw := &weighted.RRW{}
	for _, item := range c.Items {
		rrw.Add(item.Name, item.Weight)
	}
	return &WeightedRR{
		rrw: rrw,
	}
}

var _ Weight = (*WeightedRR)(nil)
