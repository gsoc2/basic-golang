package errutils

import (
	"net"
	"os"
	"syscall"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestIsIgnorableErr(t *testing.T) {
	t.Run("returns true if error is a broken pipe", func(tt *testing.T) {
		err := &net.OpError{
			Err: &os.SyscallError{
				Err: syscall.EPIPE,
			},
		}

		assert.True(tt, IsIgnorableErr(err))
	})

	t.Run("returns true if error is a connection reset error", func(tt *testing.T) {
		err := &net.OpError{
			Err: &os.SyscallError{
				Err: syscall.ECONNRESET,
			},
		}

		assert.True(tt, IsIgnorableErr(err))
	})

	t.Run("returns true if error is a timeout error", func(tt *testing.T) {
		// In production, the actual error we are getting is a non-exported
		// net.Error (http.httpError). Because the error is not exported we cannot
		// test for that specific error in our tests. DNSError is exported and we
		// can manually set the IsTimeout field so we can use this error for our
		// tests. Both DNSError and httpError are types of net.Error and adhere to
		// the timeout interface.
		err := &net.DNSError{
			IsTimeout: true,
		}

		assert.True(tt, IsIgnorableErr(err))

		wrapErr := errors.Wrap(err, "test")
		assert.True(tt, IsIgnorableErr(wrapErr))

	})

	t.Run("returns false if error is not ignorable", func(tt *testing.T) {
		err := errors.New("foo")

		assert.False(tt, IsIgnorableErr(err))
	})

	t.Run("handles wrapped errors", func(tt *testing.T) {
		err := errors.Wrap(&net.OpError{
			Err: &os.SyscallError{
				Err: syscall.EPIPE,
			},
		}, "wrapped error")

		assert.True(tt, IsIgnorableErr(err))
	})
}
