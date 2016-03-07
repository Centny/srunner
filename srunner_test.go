package srunner

import (
	"fmt"
	"github.com/Centny/gwf/routing/httptest"
	"github.com/Centny/gwf/util"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	var fcfg = util.NewFcfg3()
	fcfg.InitWithFilePath2("srunner.properties", true)
	var runner, err = NewRunner(fcfg)
	if err != nil {
		t.Error(err)
		return
	}
	runner.Start()
	time.Sleep(3 * time.Second)
	//
	fmt.Println("ab->00")
	err = runner.StopProc("SR_T3", 3000)
	if err != nil {
		t.Error(err)
		return
	}
	err = runner.StopProc("SR_T3", 3000)
	if err == nil {
		t.Error("error")
		return
	}
	//
	fmt.Println("ab->10")
	err = runner.StartProc("SR_T3", 3000)
	if err != nil {
		t.Error(err)
		return
	}
	err = runner.StartProc("SR_T3", 3000)
	if err == nil {
		t.Error("error")
		return
	}
	//
	fmt.Println("ab->20")
	err = runner.RestartProc("SR_T3", 6000)
	if err != nil {
		t.Error(err)
		return
	}
	if runner.ProcM["SR_T3"].Status != PS_RUNNING {
		t.Error("error")
		return
	}
	err = runner.StopProc("SR_T3", 6000)
	if err != nil {
		t.Error(err)
		return
	}
	//
	fmt.Println("ab->30")
	var ts = httptest.NewServer2(runner)
	res, _ := ts.G2("?exec=%v&name=%v&token=xyz", "start", "SR_T3")
	if res.IntVal("code") != 0 {
		fmt.Println(res)
		t.Error("error")
		return
	}
	res, _ = ts.G2("?exec=%v&name=%v&token=xyz", "restart", "SR_T3")
	if res.IntVal("code") != 0 {
		t.Error("error")
		return
	}
	res, _ = ts.G2("?exec=%v&name=%v&token=xyz", "stop", "SR_T3")
	if res.IntVal("code") != 0 {
		t.Error("error")
		return
	}
	//
	//test error
	err = runner.StartProc("SR_T3", 6000)
	if err != nil {
		t.Error(err)
		return
	}
	//
	fmt.Println("testing error...")
	err = runner.RestartProc("SR_T4", 6000)
	if err == nil {
		t.Error("error")
		return
	}
	err = runner.StartProc("SR_T4", 6000)
	if err == nil {
		t.Error("error")
		return
	}
	err = runner.StopProc("SR_T4", 6000)
	if err == nil {
		t.Error("error")
		return
	}
	err = runner.RestartProc("SR_T3", 100)
	if err == nil {
		t.Error("error")
		return
	}
	err = runner.StartProc("SR_xxx", 6000)
	if err == nil {
		t.Error("error")
		return
	}
	err = runner.StopProc("SR_xxx", 6000)
	if err == nil {
		t.Error("error")
		return
	}
	err = runner.RestartProc("SR_xxx", 6000)
	if err == nil {
		t.Error("error")
		return
	}
	runner.StopProc("SR_T3", 6000)
	err = runner.StartProc("SR_T3", 100)
	if err == nil {
		t.Error("error")
		return
	}
	//
	res, _ = ts.G2("?exec=%v&name=%v&token=xyz", "xxx", "SR_T3")
	if res.IntVal("code") == 0 {
		t.Error("error")
		return
	}
	res, _ = ts.G2("?exec=%v&name=%v&token=sfsd", "xxx", "SR_T3")
	if res.IntVal("code") == 0 {
		t.Error("error")
		return
	}
	res, _ = ts.G2("")
	if res.IntVal("code") == 0 {
		t.Error("error")
		return
	}

	//stop
	runner.Start()
	time.Sleep(time.Second)
	runner.Kill()
	runner.Wg.Wait()
	//
	runner.Bash = "xdds"
	runner.Start()
	time.Sleep(time.Second)
	runner.Wg.Wait()
	//
	//
	fcfg = util.NewFcfg3()
	fcfg.InitWithData(`
[SR_xx]
		`)
	_, err = NewRunner(fcfg)
	if err == nil {
		t.Error("error")
		return
	}
}
