package log

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"
)

const (
	defaultTimeLayout = "2006-01-02T15:04:05.000Z07:00"
	defaultLayout     = "%-5s [%s] %s: %s\n"
)

// ------------------------------------

type baseOutput struct {
	w io.Writer
	f Formatter
}

// NewBaseOutput ...
func NewBaseOutput(w io.Writer) Output {
	b := &baseOutput{w: w}
	return b
}

// Send ...
func (o *baseOutput) Send(e *Event) {
	if o.f != nil {
		o.w.Write([]byte(o.f.Format(e)))
		return
	}

	var msg string
	if len(e.Arguments) == 0 {
		msg = e.Format
	} else {
		msg = fmt.Sprintf(e.Format, e.Arguments...)
	}

	fmt.Fprintf(o.w, defaultLayout,
		e.Level.String(),
		e.Time.Format(defaultTimeLayout),
		e.Name,
		msg)
}

// SetFormatter ...
func (o *baseOutput) SetFormatter(f Formatter) {
	o.f = f
}

// Close ...
func (o *baseOutput) Close() {

}

// ------------------------------------

// GetQueueSize ...
func GetQueueSize(str string) int {
	mr, _ := strconv.Atoi(str)
	if mr <= 0 {
		mr = 10000
	}
	if mr > 100000 {
		mr = 100000
	}
	return mr
}

// GetBatchNum ...
func GetBatchNum(str string) int {
	mr, _ := strconv.Atoi(str)
	if mr <= 0 {
		mr = 100
	}
	if mr > 500 {
		mr = 100
	}
	return mr
}

// NewAynscOutput ...
func NewAynscOutput(w io.Writer, queueSize int, batchNum int) Output {
	o := &asyncOutput{
		evtChan:  make(chan *Event, queueSize),
		batchNum: batchNum,
	}
	o.baseOutput = &baseOutput{w: w}
	return o
}

type asyncOutput struct {
	*baseOutput
	evtChan  chan *Event
	batchNum int
	currNum  int
	buf      bytes.Buffer
	wait     sync.WaitGroup
}

func (o *asyncOutput) Send(e *Event) {
	o.evtChan <- e
}

func (o *asyncOutput) Close() {
	o.evtChan <- nil
	o.wait.Wait()
}

func (o *asyncOutput) flush() {
	bs := o.buf.Bytes()
	o.w.Write(bs)
	o.buf.Truncate(0)
	o.currNum = 0
}

func (o *asyncOutput) loop() {
	o.wait.Add(1)
	tick := time.NewTicker(2 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			o.flush()
		case evt := <-o.evtChan:
			if evt == nil {
				break
			}
			o.buf.Write(o.f.Format(evt))
			o.currNum++
			if o.currNum >= o.batchNum {
				o.flush()
			}
		}
	}
	o.wait.Done()
}
