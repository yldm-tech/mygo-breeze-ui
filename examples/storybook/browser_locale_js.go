//go:build js && wasm

package main

import "syscall/js"

func browserLocale() string {
	value := js.Global().Get("__MYGO_STORYBOOK_LOCALE__")
	if !value.Truthy() || value.Type() != js.TypeString {
		return ""
	}
	return value.String()
}
