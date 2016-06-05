package httpserver

import (
	"fmt"
	"net"
	"time"
)

// net/http/tcpKeepAliveListener
type keepAlivelListener struct {
	*net.TCPListener
}

func newKeepAliveListener(port int) (*keepAlivelListener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	tcpL, ok := l.(*net.TCPListener)
	if !ok {
		return nil, fmt.Errorf("cannot convert to tcp listener")
	}

	kaLis := &keepAlivelListener{
		TCPListener: tcpL,
	}

	return kaLis, nil
}

//
func (kl *keepAlivelListener) Accept() (c net.Conn, err error) {
	tc, err := kl.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
