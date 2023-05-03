package godom

import (
	"syscall/js"
	"time"
)

type Timeout struct {
	val js.Value
	fn  js.Func
}

func SetTimeout(fn func(), dur time.Duration) *Timeout {
	jsfn := js.FuncOf(func(_ js.Value, _ []js.Value) any {
		go fn()
		return nil
	})
	return &Timeout{
		val: global.Call("setTimeout", jsfn, int(dur/time.Millisecond)),
		fn:  jsfn,
	}
}

func (i *Timeout) Clear() {
	global.Call("clearTimeout", i.val.Int())
	i.fn.Release()
}
