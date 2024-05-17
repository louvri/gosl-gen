package process

import (
	"fmt"
	"testing"
)

func TestIsInitated(t *testing.T) {
	runner := New()
	err := runner.IsInitiated()
	if err != nil {
		fmt.Println(err)
	}
}

func TestGenerate(t *testing.T) {
	runner := New()
	err := runner.Generate("../../build/test.json")
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitialize(t *testing.T) {
	runner := New()
	err := runner.Initialize("../../build/test.json")
	if err != nil {
		t.Fatal(err)
	}
}
