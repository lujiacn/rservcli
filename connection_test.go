package rservcli

import (
	"fmt"
	"strings"
	"testing"
)

func TestLoadingR(t *testing.T) {
	r, _ := NewRcli("127.0.0.1", 6311)
	// read csv file to string
	//r.VoidExec("a<- '123'")
	//out, _ := r.Eval("a")
	//fmt.Println(out)

	// test iterator
	script := `name = c("Bob", "Mary", "Jack", "Jane")
	people = data.frame(name, ages = c(17, 23, 41, 19))
	dataframe_output = people
	if (exists("dataframe_output")) {
		if(!require(iterators)){
			install.packages("iterators")
		}
		library(iterators)
		output_iter=iter(dataframe_output, by="row")
	}

	try(nextElem(output_iter))`
	out, err := r.Eval(script)
	fmt.Println(out, err)

	out, err = r.Eval("try(nextElem(output_iter))")
	fmt.Println(out, err)

	out, err = r.Eval("try(nextElem(output_iter))")
	fmt.Println(out, err)

	out, err = r.Eval("try(nextElem(output_iter))")
	fmt.Println(out, err)

	out, err = r.Eval("try(nextElem(output_iter))")
	fmt.Printf("-%s-", out)
	if strings.Contains(out, "Error : StopIteration") {
		fmt.Println("end")
	}
}
