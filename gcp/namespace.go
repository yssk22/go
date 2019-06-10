package gcp

import "context"

var namespaceContextKey = struct{}{}

// WithNamespace sets the namespace for the current context
// and returns the new context under that namespace
func WithNamespace(ctx context.Context, namespace string) context.Context {
	return context.WithValue(ctx, namespaceContextKey, namespace)
}

// CurrentNamespace returns the namespace for the current context
func CurrentNamespace(ctx context.Context) string {
	val := ctx.Value(namespaceContextKey)
	if val != nil {
		return val.(string)
	}
	return ""
}
