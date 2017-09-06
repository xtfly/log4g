package internal

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xtfly/log4g/api"
)

func TestAsyncOutput(t *testing.T) {
	var buf bytes.Buffer
	aop := NewAsyncOutput(&buf, api.Info, 1, 1)
	f, _ := NewTextFormatter(api.CfgFormat{"type": "text", "name": "f1", "layout": "%{module}|%{lvl} >> %{msg}"})
	aop.SetFormatter(f)
	aop.Send(&api.Event{
		Format: "debug",
		Name:   "test",
		Level:  api.Debug,
		Ctx:    context.Background(),
	})
	aop.Send(&api.Event{
		Format: "abcdef",
		Name:   "test",
		Level:  api.Info,
		Ctx:    context.Background(),
	})
	time.Sleep(6 * time.Second)
	assert.Equal(t, buf.String(), "test|INF >> abcdef")
	aop.Close()
}
