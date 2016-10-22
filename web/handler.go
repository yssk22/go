package web

import "github.com/speedland/go/web/response"

type Pipeline interface {
	Next(req *Request) Response
}

// NextHandler is an alias to call the next handler in pipeline
type NextHandler func(*Request) Response

// Handler is an interface to process the request and make a response
type Handler interface {
	// Process serve http request and return the new http request and/or response value,
	Process(*Request, NextHandler) Response
}

// HandlerFunc is a func to implement Handler interface.
type HandlerFunc func(*Request, NextHandler) Response

// Process implements Handler.Process
func (h HandlerFunc) Process(r *Request, next NextHandler) Response {
	return h(r, next)
}

type handlerPipeline struct {
	head *handlerPipelineItem
	tail *handlerPipelineItem
}

func (p *handlerPipeline) Process(req *Request, next NextHandler) Response {
	// next is the next handlerPipeline to pipe
	if p.head != nil {
		return p.head.Process(req, next)
	}
	if next != nil {
		return next(req)
	}
	return nil
}

func (p *handlerPipeline) Append(handlers ...Handler) {
	for _, h := range handlers {
		item := &handlerPipelineItem{
			Node: h,
			Next: nil,
		}
		if p.head == nil {
			p.head = item
			p.tail = item
		} else {
			p.tail.Next = item
			p.tail = item
		}
	}
}

type handlerPipelineItem struct {
	Node Handler
	Next Handler
}

func (pi *handlerPipelineItem) Process(req *Request, next NextHandler) Response {
	return pi.Node.Process(req, func(r *Request) Response {
		if pi.Next == nil {
			return next(r)
		}
		return pi.Next.Process(r, next)
	})
}

// NotFoundHandler is a handler to response not found response
var NotFoundHandler = HandlerFunc(func(*Request, NextHandler) Response {
	return response.NewTextWithCode("not found", response.HTTPStatusNotFound)
})

var notFoundHandler = HandlerFunc(func(req *Request, next NextHandler) Response {
	return NotFoundHandler(req, next)
})
