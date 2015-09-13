package rservcli

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	r, _ := NewRcli("127.0.0.1", 6311)
	// fmt.Println(r)
	out, err := r.Eval(`"hello"`)
	fmt.Println(out)
	r.VoidEval(`x <- 123`)
	x, err := r.Eval(`x`)
	fmt.Println(x)
	err = r.VoidEval(`y`)
	fmt.Println(err)
	r.Close()
	out, err = r.Eval(`"test"`)
	fmt.Println("after close:", out)

}
