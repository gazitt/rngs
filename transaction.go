package main

import (
	"fmt"
	"io"
	"os"
)

type ExecFunc func(string, string) error

type Names struct {
	oldname string
	newname string
}

type Transaction struct {
	data   []Names
	fn     ExecFunc
	output io.Writer
	// mu   sync.Mutex
}

func (t *Transaction) Execute(oldname, newname string) error {
	if err := t.fn(oldname, newname); err != nil {
		return err
	}
	t.data = append(t.data, Names{oldname: oldname, newname: newname})
	return nil
}

func (t *Transaction) Rollback() {
	if t.output == nil {
		t.output = os.Stderr
	}
	for i := len(t.data) - 1; i >= 0; i-- {
		v := t.data[i]
		if err := t.fn(v.newname, v.oldname); err != nil {
			fmt.Fprintf(t.output, "Failed to Rollback: %s: %s\n", v.newname, err)
		}
	}
}
