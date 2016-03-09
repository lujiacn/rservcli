Learning from https://github.com/senseyeio/roger    

Sample:
-----------------
<code>
 r, _ := NewRcli("127.0.0.1", 6311)
 // read csv file to string
 data, _ := ioutil.ReadFile("pt_list.csv")
 // fmt.Println(result.string())
 result := string(data)
 fmt.Println(result)
 err := r.Assign("test", result)

 r.VoidEval(`out <- read.csv(text=test, header =TRUE, colClasses=c("character", "character", "character"))`)
 obj, err := r.Eval("out")
 out := obj.(map[string]interface{})
 fmt.Println(out["RECEIVED_DCM_ID"], err)
</code>

