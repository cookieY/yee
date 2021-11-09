package logger

import (
	"os"
	"testing"
)

func TestColor_Black(t *testing.T) {
	c := New()
	c.Enable()
	_, _ = os.Stdout.Write(append([]byte(c.Blue("bule")), '\n'))
}
