package liberror

import "testing"
import "fmt"

func TestInfo(t *testing.T) {
	e := New(fmt.Errorf("test error"), Generic, 10)
	if e.Info != 10 {
		t.Error("Expected 10, got ", e.Info)
	}
	e = New(fmt.Errorf("test error"), Generic, true)
	if e.Info != true {
		t.Error("Expected true, got ", e.Info)
	}
}

func createError() Error {
	return New(fmt.Errorf("test error"), NotFound, nil)
}

func TestNestedFunctionCall(t *testing.T) {
	e := createError()
	fmt.Println(e)
}

func TestType(t *testing.T) {
	e := New(fmt.Errorf("test error"), NotFound, nil)
	if e.ErrorType.String() != "NotFound" {
		t.Error("Expected NotFound, got ", e.ErrorType.String())
	}
	if e.ErrorType != NotFound {
		t.Error("Expected error type not found, got ", e.ErrorType)
	}
}

func TestRegisterError(t *testing.T) {
	e := New(fmt.Errorf("test error"), NotFound, nil)
	if e.ErrorType.String() != "NotFound" {
		t.Error("Expected NotFound, got ", e.ErrorType.String())
	}
	if e.ErrorType != NotFound {
		t.Error("Expected error type not found, got ", e.ErrorType)
	}
}
