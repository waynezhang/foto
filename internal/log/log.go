package log

import (
	"fmt"
	"sync"

	"github.com/gookit/color"
)

func Debug(format string, a ...interface{})   { Shared().debug(format, a...) }
func Info(format string, a ...interface{})    { Shared().info(format, a...) }
func Println(format string, a ...interface{}) { Shared().println(format, a...) }
func Fatal(format string, a ...interface{})   { Shared().fatal(format, a...) }

func SetVerbose(v bool) { Shared().verbose = v }

type logger struct {
	verbose bool
}

var (
	once     sync.Once
	instance logger
)

func Shared() *logger {
	once.Do(func() {
		instance = logger{verbose: false}
	})

	return &instance
}

func (l *logger) debug(format string, a ...interface{}) {
	if l.verbose {
		fmt.Printf(format+"\n", a...)
	}
}

func (l *logger) info(format string, a ...interface{}) {
	color.Info.Printf(format+"\n", a...)
}

func (l *logger) println(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func (l *logger) fatal(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	color.Println("<fg=red>ERROR</>: " + str)
}
