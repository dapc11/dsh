package framework

import (
	"time"

	"dsh/internal/terminal"
)

// InteractionRecorder records user input sequences and timing
type InteractionRecorder struct {
	keystrokes []terminal.KeyEvent
	outputs    []string
	timings    []time.Duration
	startTime  time.Time
}

// NewInteractionRecorder creates a new interaction recorder
func NewInteractionRecorder() *InteractionRecorder {
	return &InteractionRecorder{
		keystrokes: make([]terminal.KeyEvent, 0),
		outputs:    make([]string, 0),
		timings:    make([]time.Duration, 0),
		startTime:  time.Now(),
	}
}

// RecordKey records a keystroke with timing
func (r *InteractionRecorder) RecordKey(key terminal.KeyEvent) {
	r.keystrokes = append(r.keystrokes, key)
	r.timings = append(r.timings, time.Since(r.startTime))
}

// RecordOutput records terminal output
func (r *InteractionRecorder) RecordOutput(output string) {
	r.outputs = append(r.outputs, output)
}

// GetKeystrokes returns recorded keystrokes
func (r *InteractionRecorder) GetKeystrokes() []terminal.KeyEvent {
	return r.keystrokes
}

// GetOutputs returns recorded outputs
func (r *InteractionRecorder) GetOutputs() []string {
	return r.outputs
}

// Clear clears all recorded data
func (r *InteractionRecorder) Clear() {
	r.keystrokes = r.keystrokes[:0]
	r.outputs = r.outputs[:0]
	r.timings = r.timings[:0]
	r.startTime = time.Now()
}
