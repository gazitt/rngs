package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

type Mock struct {
	actual []string
}

func (m *Mock) execfunc(oldname, newname string) error {
	for i := 0; i < len(m.actual); i++ {
		if m.actual[i] == oldname {
			m.actual[i] = newname
			return nil
		}
	}
	return fmt.Errorf("Not have oldname %s", oldname)
}

func TestTransaction(t *testing.T) {
	expect := []string{"a", "b", "c", "d", "e"}
	m := Mock{
		actual: expect,
	}

	tx := Transaction{
		data:   []Names{},
		fn:     m.execfunc,
		output: ioutil.Discard,
	}

	newnames := []string{"f", "g", "h", "i", "j"}
	for i, newname := range newnames {
		if err := tx.Execute(expect[i], newname); err != nil {
			t.Errorf("Error: Execute: %s\n", err)
		}
	}

	if !reflect.DeepEqual(m.actual, newnames) {
		t.Errorf("\ngot : %s, want: %s\n", m.actual, newnames)
	}

	tx.Rollback()

	if !reflect.DeepEqual(m.actual, expect) {
		t.Errorf("\ngot : %s, want: %s\n", m.actual, expect)
	}
}
