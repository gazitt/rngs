package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
)

const (
	VAR_INDEX    = "$index"
	VAR_PATH     = "$path"
	VAR_ISDIR    = "$isdir"
	VAR_ABSPATH  = "$abspath"
	VAR_DIRNAME  = "$dirname"
	VAR_BASENAME = "$basename"
	VAR_EXTNAME  = "$extname"
)

type JSHandler struct {
	vm *otto.Otto
	js *otto.Script
}

func NewJSHandler(name, script string) (*JSHandler, error) {
	var err error
	h := &JSHandler{
		vm: otto.New(),
	}

	// If like below define a func an error occurs
	// type Func func(otto.FunctionCall) otto.Value
	// TypeError: can't convert from "string" to "otto.FunctionCall"
	// So, the type of map must be defined as follows
	// map[string]func(otto.FunctionCall) otto.Value
	m := make(map[string]func(otto.FunctionCall) otto.Value, len(funcmap))
	for k, v := range funcmap {
		m[k] = v.fn
	}
	h.vm.Set("$func", m)

	name = strings.TrimSpace(name)

	var src interface{}
	if len(script) > 0 {
		src = script
	}

	// src takes priority over file(name)
	h.js, err = h.vm.Compile(name, src)
	if err != nil {
		return h, err
	}

	return h, nil
}

func (h *JSHandler) Run(index int, oldname string) (otto.Value, error) {
	abspath, _ := filepath.Abs(oldname)
	basename := filepath.Base(oldname)
	dirname := filepath.Dir(abspath)
	extname := filepath.Ext(basename)
	info, err := os.Stat(abspath)

	h.vm.Set(VAR_INDEX, index)
	h.vm.Set(VAR_PATH, oldname)
	h.vm.Set(VAR_ABSPATH, abspath)
	h.vm.Set(VAR_DIRNAME, dirname)
	h.vm.Set(VAR_BASENAME, basename)
	h.vm.Set(VAR_EXTNAME, extname)
	h.vm.Set(VAR_ISDIR, !os.IsNotExist(err) && info.IsDir())

	return h.vm.Run(h.js)
}
