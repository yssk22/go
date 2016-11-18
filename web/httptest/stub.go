package httptest

import (
	"fmt"
	"net"
	"net/http"

	"github.com/speedland/go/x/xerrors"
	"github.com/speedland/go/x/xnet/xhttp/xhttptest"
)

// StubServer is an http server that serve the request as stub.
type StubServer struct {
	addr net.Addr
}

// Client enforce http.Client to request to the stub server
// instead of requesting external resources. mapping should be the map from external urls to stub server paths.
func (s *StubServer) Client(mapping map[string]string, c *http.Client) *http.Client {
	stubMapping := make(map[string]string)
	for k, v := range mapping {
		stubMapping[k] = fmt.Sprintf("http://%s%s", s.addr.String(), v)
	}
	return xhttptest.Stub(stubMapping, c)
}

// UseStubServer launches a stub server configured by handler
// and execute test function f.
func UseStubServer(handler http.Handler, f func(*StubServer)) {
	server := &http.Server{}
	listener, err := net.Listen("tcp", "localhost:0")
	xerrors.MustNil(err)
	defer func() {
		listener.Close()
	}()
	server.Handler = handler
	go server.Serve(listener)
	f(&StubServer{
		addr: listener.Addr(),
	})
}
