package web

import "testing"

func Test_contextKey_String(t *testing.T) {
	tests := []struct {
		name string
		c    *contextKey
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := tt.c.String(); got != tt.want {
			t.Errorf("%q. contextKey.String() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
