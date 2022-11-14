package quay

import "sync"

// Quay is a quay with scalable wharves.
type Quay[T any] struct {
	gate    sync.Mutex
	wharves map[any]*Wharf[T]
}

// New creates a virtual scalable quay.
func New[T any]() *Quay[T] {
	return &Quay[T]{}
}

// Enter enters a specific wharf, holding a ticket.
// A captain gopher is responsible to dispose the wharf if it is no longer
// needed, using the returned function.
func (q *Quay[T]) Enter(wharfID any, ticket T) (*Wharf[T], <-chan Status, func()) {
	q.gate.Lock()
	defer q.gate.Unlock()

	wharf, ok := q.wharves[wharfID]
	if !ok {
		wharf = NewWharf[T]()
		if q.wharves == nil {
			q.wharves = map[any]*Wharf[T]{
				wharfID: wharf,
			}
		} else {
			q.wharves[wharfID] = wharf
		}
	}

	return wharf, wharf.Enter(ticket), func() {
		q.gate.Lock()
		defer q.gate.Unlock()
		if wharf.idle() {
			delete(q.wharves, wharfID)
		}
	}
}
