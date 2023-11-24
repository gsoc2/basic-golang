package errutils

import (
	"io"
	"net"
	"os"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

type unwrapper interface {
	Unwrap() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Unwrap takes in an error, goes through the error cause chain, and returns the
// deepest error that has a stacktrace. This is needed because by default, every
// subsequent call to errors.Wrap overwrites the previous stacktrace. So in the
// end, we don't know where the error originated from. By going to the deepest
// one, we can find exactly where it started.
func Unwrap(err error) error {
	for {
		var u unwrapper
		if !errors.As(err, &u) {
			break
		}
		unwrapped := u.Unwrap()
		var st stackTracer
		if !errors.As(unwrapped, &st) {
			break
		}
		err = unwrapped
	}
	return err
}

// IsIgnorableErr returns true if the provided error is:
// - an EPIPE error
// - a connection reset (ECONNRESET) error
// - an http2 GOAWAY error
// - an http2 stream internal error
// - a DNS cancellation
// - a network timeout
// - an unexpected EOF
// - or a normal EOF
// IMPORTANT: When adding conditions to this function, make sure all of the
// checks are in the if and not in the return. This way, it'll try all of the
// possible conditions. Otherwise, you might get a case where it passes the
// checks in the if statement, but then fails inside of the if block, but it
// would've passed in a later condition.
func IsIgnorableErr(err error) bool {
	var serr *os.SyscallError
	if errors.As(err, &serr) && (serr.Err.Error() == syscall.EPIPE.Error() || serr.Err.Error() == syscall.ECONNRESET.Error()) {
		// EPIPE or connection reset by peer
		return true
	}

	if strings.Contains(err.Error(), "http2: server sent GOAWAY") {
		// http2 GOAWAY error: https://github.com/golang/go/issues/28930
		return true
	}

	if strings.Contains(err.Error(), "stream error: stream ID") && strings.Contains(err.Error(), "INTERNAL_ERROR") {
		// http2 stream error
		return true
	}

	var dnserr *net.DNSError
	if errors.As(err, &dnserr) && strings.Contains(dnserr.Error(), "operation was canceled") {
		// DNS cancellation
		return true
	}

	var nerr net.Error
	if errors.As(err, &nerr) && nerr.Timeout() {
		// network timeout
		return true
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		// unexpected EOF
		return true
	}

	if errors.Is(err, io.EOF) {
		// normal EOF
		return true
	}

	return false
}
