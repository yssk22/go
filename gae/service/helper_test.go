package service

import (
	"context"
	"testing"

	"github.com/yssk22/go/web/httptest"
)

func TestNewHTTPClient(t *testing.T) {
	a := httptest.NewAssert(t)
	a.NotNil(NewHTTPClient(context.Background()))
}
