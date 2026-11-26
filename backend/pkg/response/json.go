package response

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type ctxKey string

const requestIDCtxKey ctxKey = "request_id"

// WithRequestID returns a derived context carrying the request ID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDCtxKey, requestID)
}

// RequestIDFromContext retrieves the request ID stored in ctx, if any.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(requestIDCtxKey).(string)
	return id
}

// Envelope is the standardized JSON API response contract.
type Envelope struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    Meta       `json:"meta"`
}

// ErrorInfo provides machine-readable error code and human-readable message.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta contains trace and timing metadata for every response.
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

// WriteJSON serializes payload as JSON and writes it with the given HTTP status.
func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// OK sends a 200 JSON response wrapped in Envelope.
func OK(w http.ResponseWriter, r *http.Request, data any) {
	WriteJSON(w, http.StatusOK, Envelope{
		Success: true,
		Data:    data,
		Meta: Meta{
			RequestID: RequestIDFromContext(r.Context()),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// Created sends a 201 JSON response wrapped in Envelope.
func Created(w http.ResponseWriter, r *http.Request, data any) {
	WriteJSON(w, http.StatusCreated, Envelope{
		Success: true,
		Data:    data,
		Meta: Meta{
			RequestID: RequestIDFromContext(r.Context()),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// Error sends a structured error JSON response wrapped in Envelope.
func Error(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	WriteJSON(w, status, Envelope{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
		Meta: Meta{
			RequestID: RequestIDFromContext(r.Context()),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	})
}
