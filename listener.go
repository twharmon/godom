package godom

import "syscall/js"

type Listener struct {
	ty   string
	fn   js.Func
	elem *Elem
}

func (l *Listener) Remove() {
	delete(l.elem.listeners, l)
	l.elem.val.Call("removeEventListener", l.ty, l.fn)
}
