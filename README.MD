Learning from https://github.com/senseyeio/roger

Sample code
-----------------
```r
# connect to Rserve
r, err := NewRcli("127.0.0.1", 6311)
if err != nil {
    panic(fmt.Printf("Connection err: %v", err))
}

# execute R script with error message return
rScript := `test_message <- "hellow world"`
err := r.VoidExec(rScript)
if err != nil {
    fmt.Printf("Error in r script: %v", err)
}

# read R values
obj, _ := r.Eval("test_message")
fmt.Println(obj)

# obj type assertion
value, ok := obj.(string)
if !ok {
    fmt.Println("It's not ok for type string")
} else {
    fmt.Println("The value is ", value)
}

```
