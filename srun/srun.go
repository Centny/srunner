package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/smartio"
	"github.com/Centny/gwf/util"
	"github.com/Centny/srunner"
	"net/http"
	"os"
	"sync"
	"time"
)

var ef = os.Exit

func main() {
	if len(os.Args) > 2 {
		if os.Args[1] == "-c" {
			var res, err = util.HGet2("%v", os.Args[2])
			if err != nil {
				fmt.Println(err)
				ef(1)
				return
			}
			var code = res.IntVal("code")
			if code == 0 {
				fmt.Println("OK")
			} else {
				fmt.Println(util.S2Json(res))
			}
			ef(int(code))
			return
		}
	}
	var conf = "conf/srun.properties"
	if len(os.Args) > 1 {
		conf = os.Args[1]
	}
	var fcfg = util.NewFcfg3()
	fcfg.InitWithFilePath2(conf, true)
	fcfg.Print()
	redirect_l(fcfg)
	smartio.ShowLog = fcfg.Val2("showlog", "1") == "1"
	var runner, err = srunner.NewRunner(fcfg)
	if err != nil {
		fmt.Println(err)
		ef(1)
		return
	}
	runner.Start()
	var wg = sync.WaitGroup{}
	wg.Add(1)
	var listen = fcfg.Val2("listen", "")
	routing.H("^/exec(\\?.*)?", runner)
	routing.HFunc("^/_exit_$", func(hs *routing.HTTPSession) routing.HResult {
		runner.Kill()
		runner.Wg.Wait()
		wg.Done()
		log.I("srun receive exit require, the srun service will exit")
		return hs.MsgRes("OK")
	})
	var www = fcfg.Val2("www", "")
	if len(www) > 0 {
		routing.Shared.Handler("^.*$", http.FileServer(http.Dir(www)))
	}
	go func() {
		log.I("listen web on %v", listen)
		routing.ListenAndServe(listen)
	}()
	wg.Wait()
	log.I("srun done...")
	time.Sleep(time.Second)
}

func redirect_l(fcfg *util.Fcfg) {
	var out_l = fcfg.Val2("out_l", "")
	var err_l = fcfg.Val2("err_l", "")
	fmt.Printf("redirect stdout to file(%v) and stderr to file(%v)\n", out_l, err_l)
	if len(out_l) > 0 {
		smartio.RedirectStdout3(out_l)
	}
	if len(err_l) > 0 {
		smartio.RedirectStderr3(err_l)
	}
	log.SetWriter(os.Stdout)
}
