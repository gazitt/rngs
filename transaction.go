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
	exec   ExecFunc
	output io.Writer
	// When enabled Simurate NoRevert or Force option,
	// does not append and rollback.
	rollbackEnabled bool

	// mu   sync.Mutex
}

func NewTransaction(capacity int, disabled bool) *Transaction {
	t := &Transaction{
		rollbackEnabled: !disabled,
	}
	if t.rollbackEnabled {
		t.data = make([]Names, 0, capacity)
	}
	return t
}

func (t *Transaction) Execute(oldname, newname string) error {
	if err := t.exec(oldname, newname); err != nil {
		return err
	}
	if t.rollbackEnabled {
		t.data = append(t.data, Names{oldname: oldname, newname: newname})
	}
	return nil
}

func (t *Transaction) Rollback() {
	if !t.rollbackEnabled {
		return
	}

	if t.output == nil {
		t.output = os.Stderr
	}
	for i := len(t.data) - 1; i >= 0; i-- {
		v := t.data[i]
		if err := t.exec(v.newname, v.oldname); err != nil {
			fmt.Fprintf(t.output, "Failed to Rollback: %s: %s\n", v.newname, err)
		}
	}
}
