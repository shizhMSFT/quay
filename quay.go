package quay

import "sync"

// Quay is a quay with scalable wharves.
type Quay struct {
	gate    sync.Mutex
	wharves map[any]*Wharf
}

// New creates a virtual scalable quay.
func New() *Quay {
	return &Quay{}
}

// Enter enters a specific wharf, holding a ticket.
// A captain gopher is responsible to dispose the wharf if it is no longer
// needed, using the returned function.
func (q *Quay) Enter(wharfID, ticket any) (*Wharf, <-chan bool, func()) {
	q.gate.Lock()
	defer q.gate.Unlock()

	wharf, ok := q.wharves[wharfID]
	if !ok {
		wharf = NewWharf()
		if q.wharves == nil {
			q.wharves = map[any]*Wharf{
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
