liberror
========
Package liberror is a lightweight package for getting richer errors in Golang.

## Getting Started
~~~ go
package main

import "github.com/yunxing/liberror"
import "fmt"

var NetworkFail liberror.ErrorType

func simulateNetworkFail() liberror.Error {
    return New(fmt.Errorf("No network connection"), NetworkFail, nil)
}

func main() {
     // Register a new error type
    NetworkFail := liberror.RegisterError("NetworkFail")

    err := simulateNetworkFail()
    if err != nil {
        fmt.Println(err)
    }

    err = New(fmt.Errorf("File is corrupted"),
              liberror.DataCorruption, "/disk/a/report")
    fmt.Println(err)

    // You can also compare two errors by:
    if err.ErrorType == liberror.DataCorruption {
        // Get more info from the error
        fmt.Println(err.info)
    }
}

~~~

Some example outputs:
~~~
Error reason : Generic
Error message: test error
Stack:
	 _/Users/yunxing/go/src/github.com/yunxing/liberror.TestError
		 at /Users/yunxing/go/src/github.com/yunxing/liberror/error_test.go:6
	 testing.tRunner
		 at /usr/local/go/src/pkg/testing/testing.go:374
	 runtime.goexit
		 at /usr/local/go/src/pkg/runtime/proc.c:1394

~~~
