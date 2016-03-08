package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSrun(t *testing.T) {
	var ev int = -1
	ef = func(v int) {
		ev = v
	}
	os.Args = []string{"xx", "srun.properties"}
	go main()
	time.Sleep(time.Second)
	os.Args = []string{"xx", "-c", "http://127.0.0.1:3010/exec?name=SR_T3&exec=stop&token=xyz"}
	main()
	if ev != 0 {
		t.Error("error")
		return
	}
	os.Args = []string{"xx", "-c", "http://127.0.0.1:3010/exec?name=SR_T3&exec=start&token=xyz"}
	main()
	if ev != 0 {
		t.Error("error")
		return
	}
	os.Args = []string{"xx", "-c", "http://127.0.0.1:3010/exec?name=SR_T3&exec=start&token=xyz"}
	main()
	if ev == 0 {
		t.Error("error")
		return
	}
	os.Args = []string{"xx", "-c", "http://127.0.0.1:30130/exec?name=SR_T3&exec=start&token=xyz"}
	main()
	if ev == 0 {
		t.Error("error")
		return
	}
	os.Args = []string{"xx", "srun_test.properties"}
	main()
	if ev == 0 {
		t.Error("error")
		return
	}
	fmt.Println("done...")
}
