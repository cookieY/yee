package yee

import "testing"

func TestCore_Pprof(t *testing.T) {
	c := New()
	c.Pprof()
	c.Run(":9999")
}
