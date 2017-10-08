package log

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsyncOutput(t *testing.T) {
	var buf bytes.Buffer
	aop := NewAsyncOutput(&buf, Info, 1, 1)
	f, _ := NewTextFormatter(CfgFormat{"type": "text", "name": "f1", "layout": "%{module}|%{lvl} >> %{msg}"})
	aop.SetFormatter(f)
	aop.Send(&Event{
		Format: "debug",
		Name:   "test",
		Level:  Debug,
		Ctx:    context.Background(),
	})
	aop.Send(&Event{
		Format: "abcdef",
		Name:   "test",
		Level:  Info,
		Ctx:    context.Background(),
	})
	time.Sleep(6 * time.Second)
	assert.Equal(t, buf.String(), "test|INF >> abcdef")
	aop.Close()
}
