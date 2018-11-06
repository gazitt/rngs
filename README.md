
## rngs
This command-line tool uses JavaScript syntax to rename files or directories.  

## Dependencies
* [otto](https://github.com/robertkrimen/otto) - [(LICENSE)](https://github.com/robertkrimen/otto/blob/master/LICENSE)
<br/>

## Usage

* Create a Script file in the current directory  
``$ rngs -c``
<br/>

* Output the change details without actually rename.  
``$ rngs *.txt -s``
<br/>

* Sets an initial value of the $index number. (default 0)  
``$ rngs *.txt -i 100``
<br/>

* Specify the javascript file. (if not specified, using an index.js in the current directory)  
``$ rngs *.txt -f rename.js``
<br/>

* Not use the script file.  
``$ rngs *.txt -e 'flags(1)/pattern(1)/\"replacement(1)\"' -e 'flags(2)/pattern(2)/\"replacement(2)\"' ``
<br/>
<br/>
Converted into a script like below  
```javascript
    $func.join($dirname, $basename.replace(/pattern(1)/flags(1), "replacement(1)").replace(/pattern(2)/flags(2), "replacement(2)"));
```
<br/>

* Script file

    * [Target JavaScript and limitation](https://github.com/robertkrimen/otto#caveat-emptor)
    * [Regular Expression Incompatibility](https://github.com/robertkrimen/otto#regular-expression-incompatibility)

```javascript

// Return as new file name at end of this file.
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
//	 - env
//		Returns the value of environment variable
//	 - expand
//		Return the argument with an initial component of
//		tilde or environment variable replaced
//		by that user home directory.
//		$func.abs("~/foo/bar/baz");
//	 - abs
//		Returns an absolute representation of path
//	 - join
//		Joins all given path.
//		$func.join("foo","bar", "baz.txt");
//	 - exists
//		Return true if path refers to an existing path.
//		$func.exists("newname.txt");
//	 - format
//		Format returns a textual representation
//		of the time value formatted according to layout
//		https://golang.org/src/time/format.go
//		$func.format("2006-01-02T15:04:05Z07:00");
//	 - sprintf
//		Formats according to a format(https://golang.org/pkg/fmt/)
//		specifier and returns the resulting string
//		$func.sprintf("%03d", $index);

```
<br/>

* If you'd like to use "return"

```javascript

(function() {
    if $path.indexOf('foo') != -1 {
	    return $path.replace(/foo/, 'foo');
    } else if $path.indexOf('bar') != -1 {
	    return $path.replace(/bar/, 'bar');
    } else if $path.indexOf('baz') != -1 {
	    return $path.replace(/baz/, 'baz');
    }
})();

```
