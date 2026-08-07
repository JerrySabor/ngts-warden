package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type Envelope struct {
	OK        bool        `json:"ok"`
	Operation string      `json:"operation,omitempty"`
	Status    int         `json:"status,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	Error     *Error      `json:"error,omitempty"`
}

type Error struct {
	Kind      string `json:"kind"`
	Message   string `json:"message"`
	Status    int    `json:"status,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func Success(operation string, status int, data, meta interface{}) Envelope {
	return Envelope{OK: true, Operation: operation, Status: status, Data: data, Meta: meta}
}
func Failure(operation, kind, message string, status int, requestID string) Envelope {
	return Envelope{OK: false, Operation: operation, Error: &Error{Kind: kind, Message: message, Status: status, RequestID: requestID}}
}
func WriteJSON(w io.Writer, value interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(value)
}
func WriteText(w io.Writer, format string, args ...interface{}) error {
	_, err := fmt.Fprintf(w, format, args...)
	return err
}
