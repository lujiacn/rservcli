package rservcli

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	r, _ := NewRcli("127.0.0.1", 6311)
	// fmt.Println(r)
	out, err := r.Eval(`library(dplyr); test; "hello"`)
	fmt.Println(err)
	fmt.Printf("%s\n", reflect.TypeOf(out))
	fmt.Println(out)
	r.Close()
}
