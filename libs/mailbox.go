package libs

type Mailbox[M any] struct {
	queues [4]chan M
}

func NewMailbox[M any](capacity uint) Mailbox[M] {
	return Mailbox[M]{
		queues: [4]chan M{
			make(chan M, capacity),
			make(chan M, capacity),
			make(chan M, capacity),
			make(chan M, capacity),
		},
	}
}

func (m *Mailbox[M]) Push(msg M, priority int) bool {
	select {
	case m.queues[priority] <- msg:
		return true
	default:
		return false // backpressure
	}
}

func (m *Mailbox[M]) Pop() (M, bool) {
	for p := len(m.queues) - 1; p >= 0; p-- {
		select {
		case msg := <-m.queues[p]:
			return msg, true
		default:
		}
	}

	var zero M
	return zero, false
}

func (m *Mailbox[M]) Len() int {
	n := 0

	for _, q := range m.queues {
		n += len(q)
	}

	return n
}
