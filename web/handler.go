package web

// Handler is an interface to process the request and make a response
type Handler interface {
	Process(*Request) Response
}
