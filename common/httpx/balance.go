package httpx

import "sync/atomic"

type ClientType string

const (
	ClientTypeRow   ClientType = "row"
	ClientTypeProxy ClientType = "proxy"
)

type BalancePool struct {
	client  []ClientType
	current uint64
}

func (b *BalancePool) NextIndex() int {
	return int(atomic.AddUint64(&b.current, uint64(1)) % uint64(len(b.client)))
}

// GetNextPeer returns next active peer to take a connection
func (b *BalancePool) GetNextPeer() ClientType {
	// loop entire backends to find out an Alive backend
	next := b.NextIndex()
	l := len(b.client) + next // start from next and move a full cycle
	for i := next; i < l; i++ {
		idx := i % len(b.client) // take an index by modding with length
		// if we have an alive backend, use it and store if its not the original one

		if i != next {
			atomic.StoreUint64(&b.current, uint64(idx)) // mark the current one

			return b.client[idx]
		}
	}
	return ClientTypeRow
}
