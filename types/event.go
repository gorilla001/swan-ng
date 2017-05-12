package types

import (
	"fmt"
	"time"
)

// Event ...
type Event struct {
	ID     string    `json:"id"`
	Status string    `json:"status"`
	From   string    `json:"from"`
	Time   time.Time `json:"time"`
}

// Format ...
func (e *Event) Format() string {
	return fmt.Sprintf("event: swan\nid: %d\ndata: {%q:%q,%q:%q,%q:%q,%q:%q}\n\n",
		time.Now().UnixNano(),
		"id", e.ID,
		"status", e.Status,
		"from", e.From,
		"time", e.Time,
	)
}
