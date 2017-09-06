package internal

import (
	"bytes"

	"github.com/xtfly/log4g/api"
)

const (
	typeMemory = "memory"
)

type memoryOutput struct {
	api.Output
	buf bytes.Buffer
}

func (o *memoryOutput) String() string {
	return o.buf.String()
}

// NewMemoryOutput return a output instance that it print message to buffer
func NewMemoryOutput(_ api.CfgOutput) (api.Output, error) {
	r := &memoryOutput{}
	r.Output = NewBaseOutput(&r.buf, api.All)
	return r, nil
}
