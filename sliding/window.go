package sliding

import "sync"

type tickDescriptorGeneric struct {
	ts            int64
	messages      []SlidingWindowLogTaskChange
	msgQueue      chan SlidingWindowLogTaskChange
	messagesMutex sync.Mutex
}

func (t *tickDescriptorGeneric) Messages() []SlidingWindowLogTaskChange {
	return t.messages
}

func (t *tickDescriptorGeneric) AppendMessages(msg SlidingWindowLogTaskChange) {
	t.msgQueue <- msg
}

func (t *tickDescriptorGeneric) NumOfMessages() int {
	return len(t.messages)
}

func newTickDescriptorGeneric() *tickDescriptorGeneric {
	tickDesc := &tickDescriptorGeneric{
		messages:      []SlidingWindowLogTaskChange{},
		msgQueue:      make(chan SlidingWindowLogTaskChange),
		messagesMutex: sync.Mutex{},
	}

	go func() {
		for {
			msg := <-tickDesc.msgQueue

			tickDesc.messagesMutex.Lock()
			tickDesc.messages = append(tickDesc.messages, msg)
			tickDesc.messagesMutex.Unlock()
		}
	}()

	return tickDesc
}

type slidingWindowGeneric struct {
	tickDescriptor *tickDescriptorGeneric
	phase          Phase
}

func (w *slidingWindowGeneric) Collect(task SlidingWindowLogTaskChange) error {
	w.tickDescriptor.AppendMessages(task)
	return nil
}

func (w *slidingWindowGeneric) NewTick() (TickResult, error) {
	w.tickDescriptor.messagesMutex.Lock()
	defer w.tickDescriptor.messagesMutex.Unlock()

	oldPhase := w.phase
	oldMessages := append([]SlidingWindowLogTaskChange{}, w.tickDescriptor.messages...)
	w.phase = "runtime"
	w.tickDescriptor.messages = []SlidingWindowLogTaskChange{}

	return TickResult{
		Messages: oldMessages,
		Phase:    oldPhase}, nil
}

func (w *slidingWindowGeneric) Phase() string {
	return w.phase
}
