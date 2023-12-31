package log

import (
	"fmt"

	"github.com/gookit/color"
)

var verbose = false

func SetVerbose(v bool) {
	verbose = v
}

func Debug(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format+"\n", a...)
	}
}

func Info(format string, a ...interface{}) {
	color.Info.Printf(format+"\n", a...)
}

func Println(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func Fatal(format string, a ...interface{}) {
	color.Printf("<fg=red>ERROR</>: "+format+"\n", a...)
}
