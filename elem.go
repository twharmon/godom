package godom

import (
	"strings"
	"syscall/js"
)

type Elem struct {
	ty        string
	val       js.Value
	children  []*Elem
	listeners map[string]js.Func
	attrs     map[string]struct{}
	Done      chan struct{}
}

func (e *Elem) AppendChild(children ...*Elem) *Elem {
	for _, child := range children {
		e.val.Call("appendChild", child.val)
		e.children = append(e.children, child)
	}
	return e
}

func (e *Elem) Text(text any) *Elem {
	if len(e.children) == 1 && e.children[0].isTextNode() {
		e.children[0].val.Set("nodeValue", text)
		return e
	}
	e.Clear()
	n := CreateTextElem(text)
	e.val.Call("appendChild", n.val)
	e.children = append(e.children, n)
	return e
}

func (e *Elem) Attr(name string, value any) *Elem {
	e.val.Set(name, value)
	e.registerAttr(name)
	return e
}

func (e *Elem) Style(name string, value string) *Elem {
	e.val.Get("style").Set(name, value)
	e.registerAttr("style")
	return e
}

func (e *Elem) Classes(names ...string) *Elem {
	e.val.Set("classList", strings.Join(names, " "))
	e.registerAttr("class")
	return e
}

func (e *Elem) AddClass(name string) *Elem {
	e.val.Get("classList").Call("add", name)
	e.registerAttr("class")
	return e
}

func (e *Elem) RemoveClass(name string) *Elem {
	e.val.Get("classList").Call("remove", name)
	return e
}

func (e *Elem) ToggleClass(name string) *Elem {
	e.val.Get("classList").Call("toggle", name)
	e.registerAttr("class")
	return e
}

func (e *Elem) Clear() {
	for _, child := range e.children {
		e.RemoveChild(child)
		child.Clear()
	}
}

func (e *Elem) RemoveChild(child *Elem) {
	e.val.Call("removeChild", child.val)
	store.put(child)
	for i, ch := range e.children {
		if ch == child {
			if i < len(e.children)-1 {
				e.children[i] = e.children[len(e.children)-1]
			}
			e.children = e.children[:len(e.children)-1]
			break
		}
	}
}

func (e *Elem) ReplaceChild(newChild *Elem, oldChild *Elem) {
	e.val.Call("replaceChild", newChild.val, oldChild.val)
	for i, ch := range e.children {
		if ch == oldChild {
			ch.Clear()
			store.put(ch)
			e.children[i] = newChild
			break
		}
	}
}

func (e *Elem) Replace(with *Elem) {
	parent := e.val.Get("parentNode")
	parent.Call("replaceChild", with.val, e.val)
	e.Clear()
	store.put(e)
}

// func (e *Elem) Remove() {
// 	parent := e.val.Get("parentNode")
// 	parent.Call("removeChild", e.val)
// 	e.Clear()
// 	store.put(e)
// }

func (e *Elem) OnClick(cb func(*MouseEvent)) *Elem {
	e.setEventListener("click", js.FuncOf(func(_ js.Value, args []js.Value) any {
		go cb(newMouseEvent("click", args[0]))
		return nil
	}))
	return e
}

func (e *Elem) setEventListener(ty string, f js.Func) {
	e.removeEventListener(ty)
	if e.listeners == nil {
		e.listeners = make(map[string]js.Func)
	}
	e.listeners[ty] = f
	e.val.Call("addEventListener", ty, f)
}

func (e *Elem) removeEventListener(ty string) {
	if f, ok := e.listeners[ty]; ok {
		e.val.Call("removeEventListener", ty, f)
		f.Release()
	}
}

func (e *Elem) registerAttr(name string) {
	if e.attrs == nil {
		e.attrs = make(map[string]struct{})
	}
	e.attrs[name] = struct{}{}
}

func (e *Elem) isTextNode() bool {
	return e.ty == elemTypeText
}

func CreateTextElem(text any) *Elem {
	return store.get(elemTypeText, text)
}
