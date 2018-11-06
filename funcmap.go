package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
)

type Define struct {
	desc string
	fn   func(otto.FunctionCall) otto.Value
}

var (
	funcmap = map[string]Define{
		"env": Define{
			desc: "Returns the value of environment variable",
			fn:   _env,
		},
		"expand": Define{
			desc: `Return the argument with an initial component of
tilde or environment variable replaced
by that user home directory.
$func.abs("~/foo/bar/baz");`,
			fn: _expand,
		},
		"abs": Define{
			desc: "Returns an absolute representation of path",
			fn:   _abs,
		},
		"join": Define{
			desc: `Joins all given path.
$func.join("foo","bar", "baz.txt");`,
			fn: _join,
		},
		"exists": Define{
			desc: `Return true if path refers to an existing path.
$func.exists("newname.txt");`,
			fn: _exists,
		},
		"format": Define{
			desc: `Format returns a textual representation
of the time value formatted according to layout
https://golang.org/src/time/format.go
$func.format("2006-01-02T15:04:05Z07:00");`,
			fn: _format,
		},
		"sprintf": Define{
			desc: `Formats according to a format(https://golang.org/pkg/fmt/)
specifier and returns the resulting string
$func.sprintf("%03d", $index);`,
			fn: _sprintf,
		},
	}
)

func _env(call otto.FunctionCall) otto.Value {
	a := call.Argument(0).String()
	v, _ := otto.ToValue(os.Getenv(a))
	return v
}

func _expand(call otto.FunctionCall) otto.Value {
	a := call.Argument(0).String()
	switch a[0] {
	case '~':
		var u string
		if runtime.GOOS == "windows" {
			u = os.Getenv("USERPROFILE")
		} else {
			u = os.Getenv("HOME")
		}
		a = strings.Replace(a, "~", u, 1)
	case '$':
		i := 0
		for ; i < len(a); i++ {
			if a[i] == '/' {
				break
			}
		}
		v := a[1:i]
		env := os.Getenv(v)
		if len(env) > 0 {
			a = strings.Replace(a, "$"+v, env, 1)
		}
	}
	v, _ := otto.ToValue(os.Getenv(a))
	return v
}

func _abs(call otto.FunctionCall) otto.Value {
	abspath, _ := filepath.Abs(call.Argument(0).String())
	v, _ := otto.ToValue(abspath)
	return v
}

func _join(call otto.FunctionCall) otto.Value {
	a := make([]string, len(call.ArgumentList), len(call.ArgumentList))
	for i, v := range call.ArgumentList {
		a[i] = v.String()
	}
	v, _ := otto.ToValue(filepath.Join(a...))
	return v
}

func _exists(call otto.FunctionCall) otto.Value {
	a := call.Argument(0).String()
	v, _ := otto.ToValue(exists(a))
	return v
}

func _format(call otto.FunctionCall) otto.Value {
	a := call.Argument(0).String()
	s := time.Now().Format(a)
	v, _ := otto.ToValue(s)
	return v
}

func _sprintf(call otto.FunctionCall) otto.Value {
	if length := len(call.ArgumentList); length > 0 {
		format := call.ArgumentList[0].String()
		if length == 1 {
			s := format
			v, _ := otto.ToValue(s)
			return v
		}

		args := call.ArgumentList[1:]
		a := make([]interface{}, len(args), len(args))
		for i, v := range args {
			switch {
			case v.IsString():
				a[i], _ = v.ToString()
			case v.IsNumber():
				_, err := strconv.Atoi(v.String())
				if err != nil {
					a[i], _ = v.ToFloat()
				} else {
					a[i], _ = v.ToInteger()
				}
			case v.IsBoolean():
				a[i], _ = v.ToBoolean()
			case v.IsUndefined():
				a[i] = nil
			default:
				a[i] = v.String()
			}
		}

		s := fmt.Sprintf(format, a...)
		v, _ := otto.ToValue(s)
		return v
	}

	return otto.Value{}
}
