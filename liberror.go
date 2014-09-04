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
