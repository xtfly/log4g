package log

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultTimeLayout = "2006-01-02T15:04:05.000Z07:00"
	defaultLayout     = "%-5s [%s] %s: %s\n"

	queueMinSize       = 1000
	queueMaxSize       = 100000
	batchMinNum        = 20
	batchMaxNum        = 500
	flagClosed   int32 = 1
)

// ------------------------------------

type baseOutput struct {
	w io.Writer
	f Formatter
	t Level //threshold
}

// NewBaseOutput ...
func NewBaseOutput(w io.Writer, threshold Level) Output {
	b := &baseOutput{w: w, t: threshold}
	return b
}

// Send ...
func (o *baseOutput) Send(e *Event) {
	if e.Level < o.t {
		return
	}

	if o.f != nil {
		o.w.Write([]byte(o.f.Format(e)))
		return
	}

	fmt.Fprintf(o.w, defaultLayout,
		e.Level.String(),
		e.Time.Format(defaultTimeLayout),
		e.Name,
		e.Message())
}

// SetFormatter ...
func (o *baseOutput) SetFormatter(f Formatter) {
	o.f = f
}

func (o *baseOutput) CallerInfoFlag() int {
	if o.f != nil {
		return o.f.CallerInfoFlag()
	}
	return ciNoneFlog
}

// Close ...
func (o *baseOutput) Close() {

}

// ------------------------------------

// GetQueueSize ...
func GetQueueSize(str string) int {
	mr, _ := strconv.Atoi(str)
	if mr <= 0 {
		mr = queueMinSize
	}
	if mr > queueMaxSize {
		mr = queueMaxSize
	}
	return mr
}

// GetBatchNum ...
func GetBatchNum(str string) int {
	mr, _ := strconv.Atoi(str)
	if mr <= 0 {
		mr = batchMinNum
	}
	if mr > batchMaxNum {
		mr = batchMaxNum
	}
	return mr
}

// GetThresholdLvl ...
func GetThresholdLvl(str string) Level {
	lvl := LevelFrom(str)
	if lvl == Uninitialized {
		lvl = All
	}
	return lvl
}

// NewAsyncOutput ...
func NewAsyncOutput(w io.Writer, threshold Level, queueSize int, batchNum int) Output {
	o := &asyncOutput{
		evtChan:  make(chan *Event, queueSize),
		batchNum: batchNum,
	}
	o.baseOutput = &baseOutput{w: w, t: threshold}
	go o.loop()
	return o
}

type asyncOutput struct {
	*baseOutput
	evtChan  chan *Event
	batchNum int
	currNum  int
	buf      bytes.Buffer
	wait     sync.WaitGroup
	closed   int32
}

func (o *asyncOutput) Send(e *Event) {
	if atomic.LoadInt32(&o.closed) == flagClosed {
		return
	}
	o.evtChan <- e
}

func (o *asyncOutput) Close() {
	// support duplicate call Close method
	if atomic.LoadInt32(&o.closed) == flagClosed {
		return
	}
	o.evtChan <- nil
	o.wait.Wait()
	atomic.StoreInt32(&o.closed, flagClosed)
	close(o.evtChan)
}

func (o *asyncOutput) flush() {
	bs := o.buf.Bytes()
	o.w.Write(bs)
	o.buf.Truncate(0)
	o.currNum = 0
}

func (o *asyncOutput) loop() {
	o.wait.Add(1)
	defer o.wait.Done()

	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			o.flush()
		case evt := <-o.evtChan:
			if evt == nil {
				o.flush()
				return
			}

			if evt.Level >= o.t {
				o.buf.Write(o.f.Format(evt))
				o.currNum++
				if o.currNum >= o.batchNum {
					o.flush()
				}
			}
		}
	}
}
