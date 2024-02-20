package sliding

const (
	kTickTimeStartupSeconds = 2
	kTickTimeRuntimeSeconds = 10
)

// SlidingWindowLogTaskChange is an interface that represents a change in the logging level of a task in a sliding window log.
// It provides methods to retrieve the component ID, logger fully qualified domain name (FQDN), creation time of the change, old log level, and new log level.
type SlidingWindowLogTaskChange struct {
	ComponentId  string
	LoggerFqdn   string
	CreationTime int64
	OldLogLevel  string
	NewLogLevel  string
}

// TickDescriptor is an interface that represents a descriptor for the tick in a process.
// It provides methods to retrieve the timestamp of the tick (Ts) and the number of messages associated with the tick (NumOfMessages).
type TickDescriptor interface {
	NumOfMessages() int
	Messages() []SlidingWindowLogTaskChange
}

type Phase = string

// TickResult provides an extensible way to communicate new tick results required by processing functions.
type TickResult struct {
	Messages []SlidingWindowLogTaskChange
	Phase    Phase
}

// SlidingWindow represents a sliding window that keeps track of a fixed number of items.
type SlidingWindow interface {
	// Collect adds the given task in the list of tasks that must be submitted to the backend
	// at the end of the sliding window tick.
	Collect(task SlidingWindowLogTaskChange) error

	// NewTick instantiates a new sliding window internal tick. It returns the previous messages.
	NewTick() (TickResult, error)

	Phase() string
}

// NewSlidingWindow creates a new instance of SlidingWindow.
func NewSlidingWindow() SlidingWindow {
	return &slidingWindowGeneric{
		tickDescriptor: newTickDescriptorGeneric(),
		phase:          "startup",
	}
}
