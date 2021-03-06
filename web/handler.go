package web

import "github.com/yssk22/go/web/response"

// NextHandler is an alias to call the next handler in pipeline
type NextHandler func(*Request) *response.Response

// Handler is an interface to process the request and make a *response.Response
type Handler interface {
	// Process serve http request and return the new http request and/or *response.Response value,
	Process(*Request, NextHandler) *response.Response
}

// HandlerFunc is a func to implement Handler interface.
type HandlerFunc func(*Request, NextHandler) *response.Response

// Process implements Handler.Process
func (h HandlerFunc) Process(r *Request, next NextHandler) *response.Response {
	return h(r, next)
}

type handlerPipeline struct {
	head     *handlerPipelineItem
	tail     *handlerPipelineItem
	Handlers []Handler // keep all handlers to concat
}

func (p *handlerPipeline) Process(req *Request, next NextHandler) *response.Response {
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
	p.Handlers = append(p.Handlers, handlers...)
}

type handlerPipelineItem struct {
	Node Handler
	Next Handler
}

func (pi *handlerPipelineItem) Process(req *Request, next NextHandler) *response.Response {
	return pi.Node.Process(req, func(r *Request) *response.Response {
		if pi.Next == nil {
			// next can be nil if no more pipeline item is defined
			if next != nil {
				return next(r)
			}
			return nil
		}
		return pi.Next.Process(r, next)
	})
}
