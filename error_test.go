package liberror

import "testing"
import "fmt"

func TestInfo(t *testing.T) {
	e := New(fmt.Errorf("test error"), Generic, 10)
	if e.info != 10 {
		t.Error("Expected 10, got ", e.info)
	}
	e = New(fmt.Errorf("test error"), Generic, true)
	if e.info != true {
		t.Error("Expected true, got ", e.info)
	}
}

func TestType(t *testing.T) {
	e := New(fmt.Errorf("test error"), NotFound, nil)
	if e.errorType.String() != "NotFound" {
		t.Error("Expected NotFound, got ", e.errorType.String())
	}
	if e.errorType != NotFound {
		t.Error("Expected error type not found, got ", e.errorType)
	}
}
