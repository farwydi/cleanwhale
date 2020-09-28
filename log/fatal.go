package log

import "os"

// Fatal print to the standard logger followed by a call to os.Exit(1).
func Fatal(err error) {
	println("fatal: " + err.Error())
	os.Exit(1)
}
