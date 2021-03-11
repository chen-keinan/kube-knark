package trace

import "time"

//  Event Types

// processData is common to all events that are published. Fields tagged with
// omitempty are optional data.
type processData struct {
	StartTime time.Time `json:"start_time"`

	PPID       int      `json:"ppid"`
	ParentComm string   `json:"parent_comm,omitempty"`
	PID        int      `json:"pid"`
	UID        int      `json:"uid"`
	GID        int      `json:"gid"`
	Comm       string   `json:"comm,omitempty"`
	Exe        string   `json:"exe,omitempty"`
	Args       []string `json:"args,omitempty"`
}

// Proc
type ProcessStarted struct {
	Type string `json:"type"`
	processData
}

type ProcessExited struct {
	Type string `json:"type"`
	processData
	EndTime     time.Time     `json:"end_time"`
	RunningTime time.Duration `json:"running_time_ns"`
}

type ProcessError struct {
	Type string `json:"type"`
	processData
	ErrorCode int32 `json:"error_code"`
}
