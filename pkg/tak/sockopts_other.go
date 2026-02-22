//go:build !linux

package tak

import (
	"net"
	"time"
)

// setTCPUserTimeout is a no-op on non-Linux platforms.
// TCP_USER_TIMEOUT is a Linux-only socket option.
func setTCPUserTimeout(_ net.Conn, _ time.Duration) error {
	return nil
}

// setTCPUserTimeoutTLS is a no-op on non-Linux platforms.
func setTCPUserTimeoutTLS(_ net.Conn, _ time.Duration) error {
	return nil
}
