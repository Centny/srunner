package main

import (
	"fmt"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/smartio"
	"github.com/Centny/gwf/util"
	"github.com/Centny/srunner"
	"os"
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
	var runner, err = srunner.NewRunner(fcfg)
	if err != nil {
		fmt.Println(err)
		ef(1)
		return
	}
	redirect_l(fcfg)
	var listen = fcfg.Val2("listen", "")
	routing.H("^/exec(\\?.*)?", runner)
	log.I("listen web on %v", listen)
	routing.ListenAndServe(listen)
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
