package main

import (
	"fmt"
	"strings"
)

type Expression struct {
	data []string
}

func split(s string) ([]string, error) {
	_s := s
	i := 0
	a := make([]string, 3, 3)
	for len(_s) > 0 {
		n := strings.Index(_s, "/")
		if n == -1 {
			break
		}
		if n-1 >= 0 && _s[n-1] == '\\' {
			a[i] += _s[:n+1]
			_s = _s[n+1:]
			continue
		}

		a[i] += _s[:n]
		_s = _s[n+1:]
		if i >= 2 {
			break
		}
		i++
	}
	if i < 2 {
		return nil, fmt.Errorf("Expressoin: syntax error")
	}
	a[2] = _s

	return a, nil
}

func (e *Expression) Parse(s string) error {
	a, err := split(s)
	if err != nil {
		return err
	}
	if a[2] == "" {
		a[2] = "''"
	}
	e.data = append(e.data, fmt.Sprintf("replace(/%s/%s, %s)", a[1], a[0], a[2]))
	return nil
}

func (e *Expression) ToScript() string {
	if len(e.data) == 0 {
		return ""
	}

	script := VAR_BASENAME
	for _, v := range e.data {
		script += "." + v
	}
	return "$func.join($dirname, " + script + ");"
}
