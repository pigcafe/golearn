package trace

import (
  "io"
  "fmt"
)

type Tracer interface {
  Trace(...interface{})
}

type tracer struct {
  out io.Writer
}

func (t *tracer) Trace(a ...interface{})  {
  t.out.Write([]byte(fmt.Sprint(a...)))
  t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
  return &tracer{out: w}
}

type nilTracer struct{

}

func (t *nilTracer) Trace(a ...interface{})  {

}

// Return Tracer that can ignore to call Trace method.
func Off() Tracer {
  return &nilTracer{}
}