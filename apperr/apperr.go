// Package apperr defines a small error type the service layer uses to say
// what went wrong in a way the HTTP layer can turn into a status code.
package apperr

// Kind groups errors by how the caller should treat them.
type Kind int

const (
	KindInternal Kind = iota
	KindBadRequest
	KindUnauthorized
	KindNotFound
	KindConflict
)

// Error carries a client-safe message plus the underlying cause for logging.
type Error struct {
	Kind    Kind
	Message string // safe to return to the client
	Err     error  // underlying cause, kept for server-side logs
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap lets errors.Is/As reach the underlying cause.
func (e *Error) Unwrap() error { return e.Err }

func BadRequest(msg string) *Error   { return &Error{Kind: KindBadRequest, Message: msg} }
func Unauthorized(msg string) *Error { return &Error{Kind: KindUnauthorized, Message: msg} }
func NotFound(msg string) *Error     { return &Error{Kind: KindNotFound, Message: msg} }
func Conflict(msg string) *Error     { return &Error{Kind: KindConflict, Message: msg} }

// Internal wraps an unexpected error. The message stays generic so we never
// leak internals to the client; the real cause travels in Err for the logs.
func Internal(err error) *Error {
	return &Error{Kind: KindInternal, Message: "Something went wrong. Try again.", Err: err}
}
