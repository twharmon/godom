package godom

import (
	"syscall/js"
	"time"
)

type Interval struct {
	val js.Value
	fn  js.Func
}

func SetInterval(fn func(), dur time.Duration) *Interval {
	jsfn := js.FuncOf(func(_ js.Value, _ []js.Value) any {
		fn()
		return nil
	})
	return &Interval{
		val: global.Call("setInterval", jsfn, int(dur/time.Millisecond)),
		fn:  jsfn,
	}
}

func (i *Interval) Clear() {
	global.Call("clearInterval", i.val.Int())
	i.fn.Release()
}
