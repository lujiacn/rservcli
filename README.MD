Learning from https://github.com/senseyeio/roger

## Sample code

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

## Known issue (due to rserve?)

Cannot do assigment if value > 16 MB. If do string assiment, need to by batch.

```
// batchStrAssign do string assignment by each 16M
func batchStrAssign(rVar, rawData string, rClient *rservcli.Rcli) (err error) {
	batchSize := 16*1024*1024 - 1
	addRound := 0
	tLen := len(rawData)
	if tLen%batchSize > 0 {
		addRound = 1
	}

	batchNumber := (tLen / batchSize) + addRound
	batchRvar := []string{}
	for i := 0; i < batchNumber; i++ {
		var value string
		if i == batchNumber-1 {
			value = rawData[i*batchSize : i*batchSize+tLen%batchSize]
		} else {
			value = rawData[i*batchSize : i*batchSize+batchSize]
		}

		tRvar := fmt.Sprintf("%s_%v", rVar, i)
		err = rClient.Assign(tRvar, value)
		if err != nil {
			return err
		}
		batchRvar = append(batchRvar, tRvar)
	}

	//merg in R
	tScript := fmt.Sprintf("%s <- paste0(%s)", rVar, strings.Join(batchRvar, ", "))
	err = rClient.VoidExec(tScript)
	if err != nil {
		return err
	}
	return nil

}
```
