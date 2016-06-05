package httpserver

import (
	"fmt"
	"net"
	"time"
)

// net/http/tcpKeepAliveListener
type keepAliveListener struct {
	*net.TCPListener
}

func newKeepAliveListener(port int) (*keepAliveListener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	tcpL, ok := l.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("cannot convert to tcp listener")
	}

	kaLis := &keepAliveListener{
		TCPListener: tcpL,
	}

	return kaLis, nil
}

//
func (kl *keepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := kl.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
