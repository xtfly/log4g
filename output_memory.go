package log

import (
	"bytes"
)

const (
	typeMemory = "memory"
)

type memoryOutput struct {
	Output
	buf bytes.Buffer
}

func (o *memoryOutput) String() string {
	return o.buf.String()
}

// NewMemoryOutput return a output instance that it print message to buffer
func NewMemoryOutput(_ CfgOutput) (Output, error) {
	r := &memoryOutput{}
	r.Output = NewBaseOutput(&r.buf, All)
	return r, nil
}
