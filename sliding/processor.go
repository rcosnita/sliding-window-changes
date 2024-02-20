package sliding

import (
	"fmt"
	"time"
)

// HttpPayload defines the model that the backend service accepts.
type HttpPayload = []SlidingWindowLogTaskChange

// SlidingWindowProcess represents a process that operates on a sliding window.
// It takes a SlidingWindow as a parameter and regularly submits the messages to a backend service.
type SlidingWindowProcess struct {
	slidingWindow SlidingWindow
	logServiceUrl string
}

// submit takes all the messages currently contained in the queue for submission and sends them to the backend.
func (p *SlidingWindowProcess) submit(payload HttpPayload) error {
	// TODO(rcosnita) implement the actual code
	fmt.Printf("SlidingWindowProcess::submit: submitting %d number of log changes to service %s", len(payload), p.logServiceUrl)
	fmt.Println()
	return nil
}

func NewSlidingWindowProcess(window SlidingWindow, llmUrl string) *SlidingWindowProcess {
	return &SlidingWindowProcess{
		slidingWindow: window,
		logServiceUrl: llmUrl,
	}
}

// Process starts a processing loop that ticks at a regular interval of time.
func (p *SlidingWindowProcess) Process() {
	initialSleep := kTickTimeStartupSeconds * time.Second

	// TODO(rcosnita) add support for consuming signals for closing the processing loop.
	for {
		time.Sleep(initialSleep)
		result, err := p.slidingWindow.NewTick()
		if err != nil {
			panic(err)
		}

		httpPayload := []SlidingWindowLogTaskChange{}

		for _, msg := range result.Messages {
			// TODO(rcosnita) append the message to the http backend.
			httpPayload = append(httpPayload, msg)
		}

		fmt.Printf("Current phase: %s", result.Phase)
		fmt.Println()
		err = p.submit(httpPayload)
		if err != nil {
			// TODO(rcosnita) wire this over a dedicated error channel.
			panic(err)
		}

		initialSleep = kTickTimeRuntimeSeconds * time.Second
	}
}
