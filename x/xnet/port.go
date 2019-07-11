package xnet

import (
	"net"
	"github.com/yssk22/go/x/xerrors"
)

// GetEphemeralPort returns a ephemeral port numbel
func GetEphemeralPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, xerrors.Wrap(err, "cannot get an ephemeral port")
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
