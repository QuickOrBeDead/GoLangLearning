package hello_test

import (
	"testing"

	h "github.com/QuickOrBeDead/GoLangLearning/GoModules"
)

func TestHello(t *testing.T) {
	want := "Hello, world."
	if got := h.Hello(); got != want {
		t.Errorf("Hello() = %v, want %v", got, want)
	}
}
