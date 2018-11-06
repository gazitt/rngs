package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gazitt/flago"
)

const (
	SCRIPT_FILE     = "index.js"
	SCRIPT_TEMPLATE = `// Return as new file name at end of this file.
// If returns an empty string, skip this file

var pattern = /(foo|bar|baz)/i;
var replacement = $func.sprintf("$1-%03d", $index);
$func.join($dirname, $basename.replace(pattern, replacement));

// - docs: Available defined variables and functions
//
// $index   : A number of loop index
// $path    : Path
// $abspath : Absolute path
// $basename: Base name
// $dirname : Directory name
// $extname : File extension
// $isdir   : Whether the path is a directory
// $func    : Additional functions are available
`
)

type Options struct {
	File       string
	Expression Expression
	Verbose    bool
	Simulate   bool
	NoRevert   bool
	Force      bool
	Index      int
}

var (
	NAME    string
	VERSION string
	o       Options
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func init() {
	flago.Bool("version", -1, false, "Output version information and exit", func(flago.Value) error {
		return fmt.Errorf("%s %s", NAME, VERSION)
	})

	flago.Bool("create", 'c', false,
		fmt.Sprintf("Create %s in the current directory", SCRIPT_FILE),
		func(flago.Value) error {
			if exists(SCRIPT_FILE) {
				return fmt.Errorf("File already exists: %s", SCRIPT_FILE)
			}
			s := SCRIPT_TEMPLATE
			for k, v := range funcmap {
				s += "//\t - " + k + "\n"
				for _, vv := range strings.Split(v.desc, "\n") {
					s += "//\t\t" + vv + "\n"
				}
			}

			if err := ioutil.WriteFile(SCRIPT_FILE, []byte(s), 0664); err != nil {
				return err
			}
			return fmt.Errorf("Created: %s", SCRIPT_FILE)
		})

	flago.String("expression", 'e', "",
		`Specify a simple replace expression.
Can be Chain by specifying this option multiple.
Syntax: '[Flags]/[Expression]/\"[Replacement]\"'`,
		func(v flago.Value) error {
			return o.Expression.Parse(v.String())
		})

	flago.StringVar(&o.File, "name", 'n', "index.js", "Specify a script file", nil)
	flago.BoolVar(&o.Verbose, "verbose", 'v', false, "Verbose output", nil)
	flago.BoolVar(&o.Simulate, "simulate", 's', false, "Simulation. No rename is done", nil)
	flago.BoolVar(&o.NoRevert, "no-revert", 'r', false, "Does not revert even if an error occurs", nil)
	flago.BoolVar(&o.Force, "force", 'f', false, "", nil)

	flago.IntVar(&o.Index, "index", 'i', 0, "Initial value of index number",
		func(v flago.Value) error {
			if n := v.Get().(int); n < 0 {
				return fmt.Errorf("Must be greater than or equal zero: %s", v.String())
			}
			return nil
		})
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: recover: %v", err)
			os.Exit(1)
		}
	}()

	os.Exit(_main())
}

func _main() int {
	flago.Parse()

	if err := o.Do(flago.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}

	return 0
}
