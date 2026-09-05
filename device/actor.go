package device

import (
	"datacollector/config"
	"datacollector/libs"
	"datacollector/mesures"
	"sync/atomic"
	"time"
)

type Actor[S any] struct {
	id          int
	mailbox     chan Telemetry
	pool        *libs.Pool[S]
	snapshot    atomic.Pointer[S]
	updateState func(*S, Telemetry)
	queued      bool
}

func NewActor[S any](id int, pool *libs.Pool[S], updateState func(*S, Telemetry)) *Actor[S] {
	actor := &Actor[S]{
		id:          id,
		mailbox:     make(chan Telemetry, config.MAILBOX_SIZE),
		pool:        pool,
		queued:      false,
		updateState: updateState,
	}
	actor.snapshot.Store(new(S))

	return actor
}

func (a *Actor[S]) Update(quantum int) bool {
	for i := 0; i < quantum; i++ {
		select {
		case msg := <-a.mailbox:
			state := a.pool.Get()

			a.updateState(state, msg)

			// Publish a consistent snapshot for lock-free readers.
			oldState := a.snapshot.Swap(state)
			a.pool.Put(oldState)

			elapsed := time.Since(msg.TTL)
			mesures.Stats.Add(elapsed)
		default:
			a.queued = false
			return false
		}
	}

	return true
}

func (a *Actor[S]) Send(msg Telemetry) bool {
	select {
	case a.mailbox <- msg:
		// Sent immediately.
	default:
		// Mailbox is full.
		start := time.Now()
		a.mailbox <- msg
		mesures.Stats.AddMailboxBackpressure(time.Since(start))
	}

	if a.queued {
		return false
	}

	a.queued = true
	return true
}

func (a *Actor[S]) State() *S {
	snapshot := a.snapshot.Load()
	if snapshot == nil {
		return new(S)
	}

	return snapshot
}
