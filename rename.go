package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func fatalf(format string, args ...interface{}) bool {
	fmt.Fprintf(os.Stderr, format, args...)
	return true
}

func (o *Options) output(format string, args ...interface{}) {
	if o.Simulate || o.Verbose {
		fmt.Fprintf(os.Stdout, format, args...)
	}
}

func isSamepath(oldname, newname string) bool {
	oldname, _ = filepath.Abs(oldname)
	newname, _ = filepath.Abs(newname)
	if oldname == newname {
		return true
	}
	return false
}

func renamefunc(oldname, newname string) error {
	if exists(newname) {
		return fmt.Errorf("Conflict: %s", newname)
	}
	return os.Rename(oldname, newname)
}

func (o *Options) Do(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Not specified file or directory")
	}

	// Powershell is not expand wildcard. why?
	if runtime.GOOS == "windows" {
		a := make([]string, 0, len(args))
		for _, v := range args {
			vv, err := filepath.Glob(v)
			if err != nil {
				return fmt.Errorf("Glob: %s", err)
			}
			a = append(a, vv...)
		}
		args = a
	}

	script := o.Expression.ToScript()
	h, err := NewJSHandler(o.File, script)
	if err != nil {
		return err
	}

	failed := false
	// Since it is increment one by one before processing,
	// the initial value should be subtract
	index := o.Index - 1

	tx := Transaction{
		data:   make([]Names, 0, len(args)),
		output: os.Stderr,
	}

	if o.Simulate {
		if len(script) > 0 {
			fmt.Fprintf(os.Stdout, "Script: %s\n", script)
		}
		tx.fn = GetSimulationFunc(args, exists)
	} else {
		tx.fn = renamefunc
	}

	for _, oldname := range args {
		if !exists(oldname) {
			failed = fatalf("Fatal: Not exists: %s\n", oldname)
			break
		}
		index++

		result, err := h.Run(index, oldname)
		if err != nil {
			failed = fatalf("Fatal: Script: %s\n", err)
			break
		}

		newname := strings.TrimSpace(result.String())
		if newname == "" {
			o.output("Skip: %s\n", oldname)
			// should be subtract 1
			index--
			continue
		}

		if !isSamepath(oldname, newname) {
			if err := tx.Execute(oldname, newname); err != nil {
				failed = fatalf("Fatal: %s\n", err)
				break
			}
		}
		o.output("Rename: %s -> %s\n", oldname, newname)
	}

	if failed && !o.Simulate && !o.NoRevert {
		tx.Rollback()
	}
	return nil
}
