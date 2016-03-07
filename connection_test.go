package rservcli

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	r, _ := NewRcli("127.0.0.1", 6311)
	// fmt.Println(r)

	err := r.Assign("test", 123)
	obj, err := r.Eval("test")
	fmt.Println(obj, err)
}
