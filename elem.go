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
	listeners map[string]js.Func
	attrs     map[string]struct{}
}

func (e *Elem) AppendChild(children ...*Elem) *Elem {
	for _, child := range children {
		e.val.Call("appendChild", child.val)
		e.children[child] = struct{}{}
		child.parent = e
	}
	return e
}

func (e *Elem) Text(text any) *Elem {
	if e.isTextNode() {
		e.val.Set("nodeValue", text)
	} else {
		panic(fmt.Sprintf("can not set text on element of type %s", e.ty))
	}
	return e
	// If Text() allows setting one text node as children of other element then incorporate this:
	// if len(e.children) == 1 && e.children[0].isTextNode() {
	// 	e.children[0].val.Set("nodeValue", text)
	// 	return e
	// }
	// e.Clear()
	// n := CreateTextElem(text)
	// e.val.Call("appendChild", n.val)
	// e.children = append(e.children, n)
	// return e
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

func (e *Elem) OnClick(cb func(*MouseEvent)) *Elem {
	if cb == nil {
		e.removeEventListener("click")
		return e
	}
	e.setEventListener("click", js.FuncOf(func(_ js.Value, args []js.Value) any {
		go cb(newMouseEvent(args[0]))
		return nil
	}))
	return e
}

// func (e *Elem) AddMouseEventListener(ty string, cb func(*MouseEvent)) *listener {
// 	fn := js.FuncOf(func(_ js.Value, args []js.Value) any {
// 		go cb(newMouseEvent(args[0]))
// 		return nil
// 	})
// 	e.setEventListener(ty, fn)
// 	return &listener{
// 		ty:   ty,
// 		elem: e,
// 		fn:   fn,
// 	}
// }

func (e *Elem) OnInput(cb func(*InputEvent)) *Elem {
	if cb == nil {
		e.removeEventListener("input")
		return e
	}
	e.setEventListener("input", js.FuncOf(func(_ js.Value, args []js.Value) any {
		go cb(newInputEvent(args[0]))
		return nil
	}))
	return e
}

func (e *Elem) OnMouseMove(cb func(*MouseEvent)) *Elem {
	if cb == nil {
		e.removeEventListener("mousemove")
		return e
	}
	e.setEventListener("mousemove", js.FuncOf(func(_ js.Value, args []js.Value) any {
		go cb(newMouseEvent(args[0]))
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
