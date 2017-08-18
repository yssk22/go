package service

import (
	"context"
	"testing"

	"github.com/speedland/go/web/httptest"
)

func TestNewHTTPClient(t *testing.T) {
	a := httptest.NewAssert(t)
	a.NotNil(NewHTTPClient(context.Background()))
}
