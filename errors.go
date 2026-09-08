package dbbyte

import "fmt"

// The DBByte API returns {"error": string, "hint"?: string} on failure.
// There is no "code" field in the current API - do not rely on one.
//
// A 429 can mean two different things and they are NOT both retryable:
//   - Burst rate limit: RetryAfter is set (seconds), transient, safe to retry.
//   - Plan API request quota exceeded: RetryAfter is 0/unset, retrying does
//     nothing until the billing period resets or the plan is upgraded.
//
// Check RetryAfter rather than the message text to tell them apart.
type Error struct {
	Message    string
	Status     int
	Hint       string
	RetryAfter int // seconds; 0 means "not present"
}

func (e *Error) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("dbbyte: %s (status %d, hint: %s)", e.Message, e.Status, e.Hint)
	}
	return fmt.Sprintf("dbbyte: %s (status %d)", e.Message, e.Status)
}
