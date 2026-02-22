package tak

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/sys/unix"
)

// setTCPUserTimeout sets TCP_USER_TIMEOUT on a TCP connection.
// If unacked data sits in the send buffer longer than timeout, the kernel kills the connection.
// This is a Linux-only socket option.
func setTCPUserTimeout(conn net.Conn, timeout time.Duration) error {
	if timeout <= 0 {
		return nil
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("cannot set TCP_USER_TIMEOUT: not a *net.TCPConn (got %T)", conn)
	}

	rawConn, err := tcpConn.SyscallConn()
	if err != nil {
		return fmt.Errorf("failed to get syscall conn: %w", err)
	}

	var setsockoptErr error
	if controlErr := rawConn.Control(func(fd uintptr) {
		setsockoptErr = unix.SetsockoptInt(
			int(fd),
			unix.IPPROTO_TCP,
			unix.TCP_USER_TIMEOUT,
			int(timeout.Milliseconds()),
		)
	}); controlErr != nil {
		return fmt.Errorf("rawconn control failed: %w", controlErr)
	}

	if setsockoptErr != nil {
		return fmt.Errorf("setsockopt TCP_USER_TIMEOUT failed: %w", setsockoptErr)
	}
	return nil
}

// setTCPUserTimeoutTLS sets TCP_USER_TIMEOUT on the underlying TCP connection of a TLS connection.
func setTCPUserTimeoutTLS(conn net.Conn, timeout time.Duration) error {
	if timeout <= 0 {
		return nil
	}

	// tls.Conn wraps a net.Conn; NetConn() exposes the underlying conn (Go 1.18+)
	type netConner interface {
		NetConn() net.Conn
	}
	nc, ok := conn.(netConner)
	if !ok {
		return fmt.Errorf("cannot get underlying net.Conn from TLS connection (type %T)", conn)
	}

	return setTCPUserTimeout(nc.NetConn(), timeout)
}
