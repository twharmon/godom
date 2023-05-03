package godom

import (
	"sync"
)

type elemStore struct {
	elems map[string][]*Elem
	mu    sync.Mutex
}

func newElemStore() *elemStore {
	return &elemStore{
		elems: make(map[string][]*Elem),
	}
}

func (es *elemStore) put(e *Elem) {
	if e.isTextNode() {
		e.val.Set("nodeValue", nil)
	} else {
		e.Clear()
		e.val.Set("value", nil)
		for ty := range e.listeners {
			e.removeEventListener(ty)
		}
		for attr := range e.attrs {
			e.val.Call("removeAttribute", attr)
		}
	}
	es.mu.Lock()
	es.elems[e.ty] = append(es.elems[e.ty], e)
	es.mu.Unlock()
}

func (es *elemStore) get(ty string, texts ...any) *Elem {
	es.mu.Lock()
	defer es.mu.Unlock()
	vals := es.elems[ty]
	if len(vals) > 0 {
		var e *Elem
		e, es.elems[ty] = vals[len(vals)-1], vals[:len(vals)-1]
		if e.isTextNode() {
			e.val.Set("nodeValue", texts[0])
		}
		return e
	}
	var e Elem
	e.ty = ty
	if e.isTextNode() {
		e.val = document.Call("createTextNode", texts[0])
	} else {
		e.val = document.Call("createElement", ty)
	}
	return &e
}
