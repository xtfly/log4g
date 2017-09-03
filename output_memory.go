package log

import (
	"bytes"
)

type memoryOutput struct {
	*baseOutput
	buf bytes.Buffer
}

func (o *memoryOutput) String() string {
	return o.buf.String()
}

// NewMemoryOutput return a output instance that it print message to buffer
func NewMemoryOutput(_ *CfgOutput) (Output, error) {
	r := &memoryOutput{}
	r.baseOutput = &baseOutput{w: &r.buf}
	return r, nil
}
