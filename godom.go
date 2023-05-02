package godom

import (
	"syscall/js"
)

var global = js.Global()
var document = global.Get("document")

var store = newElemStore()

const elemTypeText string = "_TEXT_"

func Create(tag string) *Elem {
	return store.get(tag)
}

func Mount(selector string, e *Elem) {
	val := document.Call("querySelector", selector)
	c := &Elem{
		val: val,
	}
	c.AppendChild(e)
}
