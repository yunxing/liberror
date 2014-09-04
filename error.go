package liberror

import "runtime"
import "fmt"
import "bytes"

type StackTrace []*runtime.Func

var (
	errorTypeNames []string       = make([]string, 0, 64)
	in             chan string    = make(chan string)
	out            chan ErrorType = make(chan ErrorType)
)

type ErrorType uint64

func (t ErrorType) String() string {
	return errorTypeNames[t]
}

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

func RegisterError(typeName string) ErrorType {
	in <- typeName
	return <-out
}

func handleRegister() {
	for {
		name := <-in
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
	err       error
	trace     *StackTrace
	errorType ErrorType
	info      interface{}
}

func New(err error, errorType ErrorType, userInfo interface{}) Error {
	stack := make([]uintptr, 1024)
	n := runtime.Callers(2, stack)
	trace := make(StackTrace, 0, n)
	for _, pc := range stack {
		if f := runtime.FuncForPC(pc); f != nil {
			trace = append(trace, f)
		}
	}
	return Error{err, &trace, errorType, userInfo}
}

func (e Error) Type() ErrorType {
	return e.errorType
}

func (e Error) StackTrace() *StackTrace {
	return e.trace
}

func (e Error) Info() interface{} {
	return e.info
}

func (e Error) Error() string {
	bs := bytes.NewBufferString("")
	bs.WriteString(fmt.Sprintf("Error reason : %s\nError message: %s \n",
		errorTypeNames[e.errorType], e.err.Error()))
	bs.WriteString("Stack: \n")
	for _, f := range *e.trace {
		file, line := f.FileLine(f.Entry())
		bs.WriteString(fmt.Sprintf("\t %s\n\t\t at %s:%d\n",
			f.Name(), file, line))
	}
	return bs.String()
}
