package node_server

import (
	"net"
	"time"

	asl "github.com/Laboratory-for-Safe-and-Secure-Systems/go-wolfssl/asl"
)

type ASLConn struct {
	tcpConn    *net.TCPConn
	aslSession *asl.ASLSession
}

// ----------------CONN INTERFACE IMPLEMENTATION---------------------//
func (c ASLConn) Read(b []byte) (n int, err error) {
	return asl.ASLReceive(c.aslSession, b)
}
func (c ASLConn) Write(b []byte) (n int, err error) {
	err = asl.ASLSend(c.aslSession, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}
func (c ASLConn) Close() error {
	asl.ASLCloseSession(c.aslSession)
	asl.ASLFreeSession(c.aslSession)
	return c.tcpConn.Close()
}
func (c ASLConn) LocalAddr() net.Addr {
	return c.tcpConn.LocalAddr()
}

func (c ASLConn) RemoteAddr() net.Addr {
	return c.tcpConn.RemoteAddr()
}

func (c ASLConn) SetDeadline(t time.Time) error {
	return c.tcpConn.SetDeadline(t)

}

func (c ASLConn) SetReadDeadline(t time.Time) error {
	return c.tcpConn.SetDeadline(t)
}
func (c ASLConn) SetWriteDeadline(t time.Time) error {
	return c.tcpConn.SetWriteDeadline(t)
}

//----------------END CONN INTERFACE IMPLEMENTATION---------------------//

type ASLListener struct {
	L  *net.TCPListener
	Ep *asl.ASLEndpoint
}

func (l ASLListener) Accept() (net.Conn, error) {
	c, err := l.L.Accept()
	if err != nil {
		return nil, err
	}

	file, _ := c.(*net.TCPConn).File()
	fd := int(file.Fd())
	session := asl.ASLCreateSession(l.Ep, fd)
	aslConn := ASLConn{
		tcpConn:    c.(*net.TCPConn),
		aslSession: session,
	}
	return aslConn, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (listener ASLListener) Close() error {
	return listener.L.Close()
}

// Addr returns the listener's network address.
func (listener ASLListener) Addr() net.Addr {
	return listener.L.Addr()
}
