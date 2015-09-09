package rservcli

import (
	"fmt"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	r, _ := NewRcli("127.0.0.1", 6311)
	// fmt.Println(r)
	out := r.Eval(`c(1,2)`)
	v, err := out.GetResultObject()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("%s", reflect.TypeOf(v))
	fmt.Println(v)
	r.Close()
}

// func TestRepareCmdString(t *testing.T) {
// script := `test <- c(1,2,3,4)
// write.csv(test, "/home/vagrant/temp.csv")`
// // out := prepareStrCmd(script)
// cli := NewRcli("127.0.0.1", 6311)
// err := cli.Connect()
// fmt.Println("err:", err)
// cli.parseInitMsg()
// p := cli.sendCommand(2, script)
// fmt.Println(p.GetResultObject())
// fmt.Println(p)
// cli.Close()
// }
