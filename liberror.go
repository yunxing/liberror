// Package liberror is a lightweight package for providing errors with richer information in Golang.
//
// For a full guide visit http://github.com/yunxing/liberror
//
//  package main
//
//  import "github.com/yunxing/liberror"
//  import "fmt"
//
//  var NetworkFail liberror.ErrorType
//
//  func simulateNetworkFail() liberror.Error {
//      return New(fmt.Errorf("No network connection"), NetworkFail, nil)
//  }
//
//  func main() {
//      // Register a new error type
//      NetworkFail := liberror.RegisterError("NetworkFail")
//
//      err := simulateNetworkFail()
//      if err != nil {
//          fmt.Println(err)
//      }
//
//      err = New(fmt.Errorf("File is corrupted"),
//                liberror.DataCorruption, "/disk/a/report")
//      fmt.Println(err)
//
//      // You can also compare two errors by:
//      if err.ErrorType == liberror.DataCorruption {
//          // Get more info from the error
//          fmt.Println(err.info)
//      }
//  }
//
// Some sample output:
//   ~/g/s/g/y/liberror (master|…) $ go test
//   Error reason : Generic
//   Error message: test error
//   Stack:
//       _/Users/yunxing/go/src/github.com/yunxing/liberror.TestError
//           at /Users/yunxing/go/src/github.com/yunxing/liberror/error_test.go:6
//       testing.tRunner
//           at /usr/local/go/src/pkg/testing/testing.go:374
//       runtime.goexit
//           at /usr/local/go/src/pkg/runtime/proc.c:1394
//
package liberror

import "runtime"
import "fmt"
import "bytes"

var (
	errorTypeNames []string       = make([]string, 0, 64)
	in             chan string    = make(chan string)
	out            chan ErrorType = make(chan ErrorType)
)

type ErrorType uint64

func (t ErrorType) String() string {
	return errorTypeNames[t]
}

// Predefined errors. These errors are designed to be as generic
// as possible to hide implementation details.
var (
	Generic              ErrorType
	NotFound             ErrorType
	PermissionDenied     ErrorType
	AuthenticationFailed ErrorType
	DataCorruption       ErrorType
	Cancelled            ErrorType
	Expired              ErrorType
	ServiceUnavailable   ErrorType
	FileSystem           ErrorType
)

// RegisterError defines a client's own generic error type.
// It takes a string as a description of this error and
// returns a unique type id for that type.
func RegisterError(typeName string) ErrorType {
	in <- typeName
	return <-out
}

func handleRegister() {
loop:
	for {
		name := <-in
		for i, n := range errorTypeNames {
			if n == name {
				out <- ErrorType(i)
				continue loop
			}
		}
		errorTypeNames = append(errorTypeNames, name)
		out <- ErrorType(len(errorTypeNames) - 1)
	}
}

func init() {
	go handleRegister()
	Generic = RegisterError("Generic")
	NotFound = RegisterError("NotFound")
	PermissionDenied = RegisterError("PermissionDenied")
	AuthenticationFailed = RegisterError("AuthenticationFailed")
	DataCorruption = RegisterError("DataCorruption")
	Cancelled = RegisterError("Cancelled")
	Expired = RegisterError("Expired")
	ServiceUnavailable = RegisterError("ServiceUnavailable")
	FileSystem = RegisterError("FileSystem")
}

type Error struct {
	Err       error
	Trace     []*runtime.Func
	ErrorType ErrorType
	Info      interface{}
}

// Wrap an error with more information
func New(err error, errorType ErrorType, userInfo interface{}) Error {
	stack := make([]uintptr, 1024)
	n := runtime.Callers(2, stack)
	trace := make([]*runtime.Func, 0, n)
	for _, pc := range stack {
		if f := runtime.FuncForPC(pc); f != nil {
			trace = append(trace, f)
		}
	}
	return Error{err, trace, errorType, userInfo}
}

func (e Error) Error() string {
	bs := bytes.NewBufferString("")
	bs.WriteString(fmt.Sprintf("Error reason : %s\nError message: %s \n",
		errorTypeNames[e.ErrorType], e.Err.Error()))
	bs.WriteString("Stack: \n")
	for _, f := range e.Trace {
		file, line := f.FileLine(f.Entry())
		bs.WriteString(fmt.Sprintf("\t %s\n\t\t at %s:%d\n",
			f.Name(), file, line))
	}
	return bs.String()
}
