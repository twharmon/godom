package godom

import "syscall/js"

type listener struct {
	ty   string
	fn   js.Func
	elem *Elem
}

func (l *listener) Remove() {
	l.elem.val.Call("removeEventListener", l.ty, l.fn)
	// l.elem.listeners[l.ty]
}
