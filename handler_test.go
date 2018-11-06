package main

import "testing"

func TestHandlerDefineVarAndFunc(t *testing.T) {
	const script = `

var prefix = "error in script:";

if ($index != 100) {
	console.log(prefix, "\n\tgot :", $index, "\n\twant:100");
}

if ($path != "name.txt") {
	console.log(prefix, "\n\tgot :", $index, "\n\twant:100");
}

var abspath = $func.abs($path);
if ($abspath != abspath) {
	console.log(prefix, "\n\t$abspath:", $abspath, "\n\tabspath :", abspath);
}

`

	h, err := NewJSHandler("", script)
	if err != nil {
		t.Errorf("Error: NewJSHandler: %s\n", err)
	}
	if _, err := h.Run(100, "name.txt"); err != nil {
		t.Errorf("Error: Run: %s\n", err)
	}
}
