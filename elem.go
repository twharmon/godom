package godom

import (
	"fmt"
	"strings"
	"syscall/js"
)

type Elem struct {
	ty        string
	val       js.Value
	parent    *Elem
	children  map[*Elem]struct{}
	listeners map[*Listener]struct{}
	attrs     map[string]struct{}
}

func (e *Elem) AppendChild(children ...*Elem) {
	for _, child := range children {
		e.val.Call("appendChild", child.val)
		e.children[child] = struct{}{}
		child.parent = e
	}
}

func (e *Elem) Text(text any) {
	if e.isTextNode() {
		e.val.Set("nodeValue", text)
	} else {
		panic(fmt.Sprintf("can not set text on element of type %s", e.ty))
	}
}

func (e *Elem) Attr(name string, value any) {
	e.val.Call("setAttribute", name, value)
	e.registerAttr(name)
}

func (e *Elem) RemoveAttr(name string) {
	e.val.Call("removeAttribute", name)
	e.deregisterAttr(name)
}

func (e *Elem) Style(name string, value string) {
	e.val.Get("style").Set(name, value)
	e.registerAttr("style")
}

func (e *Elem) Class(name string) {
	e.val.Call("setAttribute", "class", name)
	e.registerAttr("class")
}

func (e *Elem) Classes(names ...string) {
	e.val.Set("classList", strings.Join(names, " "))
	e.registerAttr("class")
}

func (e *Elem) AddClass(name string) {
	e.val.Get("classList").Call("add", name)
	e.registerAttr("class")
}

func (e *Elem) RemoveClass(name string) {
	e.val.Get("classList").Call("remove", name)
}

func (e *Elem) ToggleClass(name string) {
	e.val.Get("classList").Call("toggle", name)
	e.registerAttr("class")
}

func (e *Elem) Clear() {
	for child := range e.children {
		e.RemoveChild(child)
		child.Clear()
	}
}

func (e *Elem) RemoveChild(child *Elem) {
	e.val.Call("removeChild", child.val)
	store.put(child)
	delete(e.children, child)
}

func (e *Elem) ReplaceWith(new *Elem) {
	e.val.Call("replaceWith", new.val)
	store.put(e)
	e.parent.replaceChild(e, new)
}

func (e *Elem) replaceChild(old *Elem, new *Elem) {
	delete(e.children, old)
	e.children[new] = struct{}{}
	new.parent = e
}

func (e *Elem) AddMouseEventListener(ty string, cb func(*MouseEvent)) *Listener {
	fn := js.FuncOf(func(_ js.Value, args []js.Value) any {
		go cb(newMouseEvent(args[0]))
		return nil
	})
	l := &Listener{
		ty:   ty,
		elem: e,
		fn:   fn,
	}
	e.addEventListener(l)
	return l
}

func (e *Elem) AddInputEventListener(ty string, cb func(*InputEvent)) *Listener {
	fn := js.FuncOf(func(_ js.Value, args []js.Value) any {
		go cb(newInputEvent(args[0]))
		return nil
	})
	l := &Listener{
		ty:   ty,
		elem: e,
		fn:   fn,
	}
	e.addEventListener(l)
	return l
}

func (e *Elem) addEventListener(l *Listener) {
	e.listeners[l] = struct{}{}
	e.val.Call("addEventListener", l.ty, l.fn)
}

func (e *Elem) registerAttr(name string) {
	if e.attrs == nil {
		e.attrs = make(map[string]struct{})
	}
	e.attrs[name] = struct{}{}
}

func (e *Elem) deregisterAttr(name string) {
	delete(e.attrs, name)
}

func (e *Elem) isTextNode() bool {
	return e.ty == elemTypeText
}

func CreateTextElem(text any) *Elem {
	return store.get(elemTypeText, text)
}
