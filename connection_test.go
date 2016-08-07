package rservcli

import (
	"fmt"
	"testing"
)

func TestLoadingR(t *testing.T) {
	r, _ := NewRcli("127.0.0.1", 6311)
	// read csv file to string
	r.VoidExec("a<- '123'")
	out, _ := r.Eval("a")
	fmt.Println(out)
}
