package godom

import (
	"syscall/js"
)

var global = js.Global()
var document = global.Get("document")
var head = document.Get("head")
var body = document.Get("body")

var Document = &Elem{
	val:       document,
	children:  make(map[*Elem]struct{}),
	listeners: make(map[*Listener]struct{}),
}

var Head = &Elem{
	val:       head,
	children:  make(map[*Elem]struct{}),
	listeners: make(map[*Listener]struct{}),
}

var Body = &Elem{
	val:       body,
	children:  make(map[*Elem]struct{}),
	listeners: make(map[*Listener]struct{}),
}

var store = newElemStore()

const elemTypeText string = "_TEXT_"

func Create(tag string) *Elem {
	return store.get(tag)
}

func Mount(selector string, e *Elem) {
	val := document.Call("querySelector", selector)
	c := &Elem{
		val:       val,
		children:  make(map[*Elem]struct{}),
		listeners: make(map[*Listener]struct{}),
	}
	c.AppendChild(e)
}
