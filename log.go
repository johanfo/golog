package log

import (
	"fmt"
	golog "log"
	"os"
	"sync"
)

// Outputer is a interface that ensures we can print with call stack
type Outputer interface {
	Output(i int, s string) error
}

var (
	// Verbose triggers printing og debug info
	Verbose     = false
	createflags = golog.LstdFlags | golog.Lshortfile
	Ilog        = CreateMultiplePrint(golog.New(os.Stdout, "I:", createflags))
	Dlog        = CreateMultiplePrint(golog.New(os.Stdout, "D:", createflags))
	Wlog        = CreateMultiplePrint(golog.New(os.Stdout, "W:", createflags))
	Flog        = CreateMultiplePrint(golog.New(os.Stdout, "C:", createflags))
)

// Bits or'ed together to control what's printed.
// There is no control over the order they appear (the order listed
// here) or the format they present (as described in the comments).
// The prefix is followed by a colon only when Llongfile or Lshortfile
// is specified.
// For example, flags Ldate | Ltime (or LstdFlags) produce,
//	2009/01/23 01:23:23 message
// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// SetFlags Recreate the outputs with new flags
func SetFlags(flags int) {
	createflags = flags
	Ilog = CreateMultiplePrint(golog.New(os.Stdout, "I:", flags))
	Dlog = CreateMultiplePrint(golog.New(os.Stdout, "D:", flags))
	Wlog = CreateMultiplePrint(golog.New(os.Stdout, "W:", flags))
	Flog = CreateMultiplePrint(golog.New(os.Stdout, "C:", flags))
}

// Reset sets all print output streams to zerovalue. Effectivly preventing any output
func Reset() {
	Ilog = &MultiplePrint{}
	Dlog = &MultiplePrint{}
	Wlog = &MultiplePrint{}
	Flog = &MultiplePrint{}
}

// AppendFileWriter writes the log to the spesified filename
func AppendFileWriter(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	Ilog.Append(golog.New(f, "I:", golog.LstdFlags|golog.Lshortfile))
	Dlog.Append(golog.New(f, "D:", golog.LstdFlags|golog.Lshortfile))
	Wlog.Append(golog.New(f, "W:", golog.LstdFlags|golog.Lshortfile))
	Flog.Append(golog.New(f, "C:", golog.LstdFlags|golog.Lshortfile))
	return nil
}

// AppendFileDescriptor writes the log to a specific file descriptor
func AppendFileDescriptor(f *os.File) {
	Ilog.Append(golog.New(f, "I:", golog.LstdFlags|golog.Lshortfile))
	Dlog.Append(golog.New(f, "D:", golog.LstdFlags|golog.Lshortfile))
	Wlog.Append(golog.New(f, "W:", golog.LstdFlags|golog.Lshortfile))
	Flog.Append(golog.New(f, "C:", golog.LstdFlags|golog.Lshortfile))
}

// MultiplePrint is an Outputer that supports stacking of multiple outputs
type MultiplePrint struct {
	outs []Outputer
}

// CreateMultiplePrint takes an Output and embeds it in an MultiplePrint
func CreateMultiplePrint(o Outputer) *MultiplePrint {
	return &MultiplePrint{outs: []Outputer{o}}
}

// Output outputs to all outs
func (d *MultiplePrint) Output(i int, s string) error {
	exclusiv.Lock()
	for _, v := range d.outs {
		v.Output(i+1, s)
	}
	exclusiv.Unlock()
	return nil
}

// Append an Outputer to the list
func (d *MultiplePrint) Append(o Outputer) {
	d.outs = append(d.outs, o)
}

var exclusiv sync.Mutex

// Info logging
func Info(x ...interface{}) {
	Ilog.Output(2, fmt.Sprint(x...))
}

// Infof logging
func Infof(format string, x ...interface{}) {
	Ilog.Output(2, fmt.Sprintf(format, x...))
}

// Debug logging
func Debug(x ...interface{}) {
	if Verbose {
		Dlog.Output(2, fmt.Sprint(x...))
	}
}

// Debugf logging
func Debugf(format string, x ...interface{}) {
	if Verbose {
		Dlog.Output(2, fmt.Sprintf(format, x...))
	}
}

// Warning logging
func Warning(x ...interface{}) {
	Wlog.Output(2, fmt.Sprint(x...))
}

// Warningf logging with formatting
func Warningf(format string, x ...interface{}) {
	Wlog.Output(2, fmt.Sprintf(format, x...))
}

// Fatal logging, with exit
func Fatal(x ...interface{}) {
	Flog.Output(2, fmt.Sprint(x...))
	os.Exit(1)
}

// Fatalf logging, with exit
func Fatalf(format string, x ...interface{}) {
	Flog.Output(2, fmt.Sprintf(format, x...))
	os.Exit(1)
}

// Println supports original "log" package style
func Println(x ...interface{}) {
	Ilog.Output(2, fmt.Sprintln(x...))
}

// Printf supports original "log" package style
func Printf(format string, x ...interface{}) {
	Ilog.Output(2, fmt.Sprintf(format, x...))
}

// PrintfLevel makes it possible to print while refering to
// code up level above the actual code.
func PrintfLevel(up int, format string, x ...interface{}) {
	Ilog.Output(2+up, fmt.Sprintf(format, x...))
}
