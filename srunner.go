package srunner

import (
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/smartio"
	"github.com/Centny/gwf/util"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	PS_NOT_START  = "NOT_START"
	PS_RUNNING    = "RUNNING"
	PS_START_ERR  = "START_ERR"
	PS_RESTARTING = "RESTARTING"
)

type Proc struct {
	Name   string    `json:"name"`
	Alias  string    `json:"alias"`
	Exec   string    `json:"exec"`
	Args   string    `json:"args"`
	Envs   string    `json:"envs"`
	Cws    string    `json:"cws"`
	OutF   string    `json:"outf"`
	ErrF   string    `json:"errf"`
	Desc   string    `json:"desc"`
	On     bool      `json:"on"`
	Cmd    *exec.Cmd `json:"-"`
	Status string    `json:"status"`
}

func ParseProcL(fcfg *util.Fcfg) ([]*Proc, map[string]*Proc, error) {
	var proc_l = []*Proc{}
	var proc_m = map[string]*Proc{}
	for _, srv := range fcfg.Seces {
		if !strings.HasPrefix(srv, "SR_") {
			continue
		}
		var proc = &Proc{}
		proc.Exec = fcfg.Val2(srv+"/exec", "")
		if len(proc.Exec) < 1 {
			return nil, nil, util.Err("Parsing process fail with %v/exec is empty", srv)
		}
		proc.Name = srv
		proc.Alias = fcfg.Val2(srv+"/alias", "")
		proc.Args = fcfg.Val2(srv+"/args", "")
		proc.Envs = fcfg.Val2(srv+"/envs", "")
		proc.OutF = fcfg.Val2(srv+"/outf", "")
		proc.ErrF = fcfg.Val2(srv+"/errf", "")
		proc.Desc = fcfg.Val2(srv+"/desc", "")
		proc.Cws = fcfg.Val2(srv+"/cws", ".")
		proc.On = fcfg.Val2(srv+"/on", "0") == "1"
		proc.Status = PS_NOT_START
		proc_l = append(proc_l, proc)
		proc_m[srv] = proc
	}
	return proc_l, proc_m, nil
}

type Runner struct {
	WDelay int64
	Token  map[string]int

	Wg    sync.WaitGroup
	ProcL []*Proc
	ProcM map[string]*Proc
}

func NewRunner(fcfg *util.Fcfg) (*Runner, error) {
	var wdelay = fcfg.Int64ValV("wdelay", 8000)
	var tokens = map[string]int{}
	for _, token := range strings.Split(fcfg.Val2("token", ""), ",") {
		tokens[token] = 1
	}
	var proc_l, proc_m, err = ParseProcL(fcfg)
	var runner *Runner = nil
	if err == nil {
		runner = &Runner{
			WDelay: wdelay,
			Token:  tokens,
			Wg:     sync.WaitGroup{},
			ProcL:  proc_l,
			ProcM:  proc_m,
		}
		log.I("creat runner with %v process", len(proc_l))
	}
	return runner, err
}

func (r *Runner) Start() {
	var started int = 0
	for _, proc := range r.ProcL {
		if proc.On && proc.Status != PS_RUNNING {
			go r.RunProc(proc)
			started++
		}
	}
	log.D("Runner call start to %v process", started)
}

func (r *Runner) RestartProc(name string, timeout int64) error {
	proc, ok := r.ProcM[name]
	if !ok {
		return util.Err("Runner stop process by name(%v) fail with name is not found", name)
	}
	if !proc.On {
		return util.Err("Runner stop process by name(%v) fail with process is not active", name)
	}
	if proc.Status == PS_RUNNING {
		err := r.StopProc(name, timeout)
		if err != nil {
			return err
		}
	}
	return r.StartProc(name, timeout)
}
func (r *Runner) StartProc(name string, timeout int64) error {
	proc, ok := r.ProcM[name]
	if !ok {
		return util.Err("Runner start process by name(%v) fail with name is not found", name)
	}
	if !proc.On {
		return util.Err("Runner start process by name(%v) fail with process is not active", name)
	}
	if proc.Status == PS_RUNNING {
		return util.Err("Runner start process by name(%v) fail with process is running", name)
	}
	go r.RunProc(proc)
	var used int64 = 0
	for proc.Status != PS_RUNNING {
		time.Sleep(100 * time.Millisecond)
		used += 100
		if used >= timeout {
			return util.Err("Runner start process(%v) timeout", name)
		}
	}
	return nil
}
func (r *Runner) StopProc(name string, timeout int64) error {
	proc, ok := r.ProcM[name]
	if !ok {
		return util.Err("Runner stop process by name(%v) fail with name is not found", name)
	}
	if !proc.On {
		return util.Err("Runner stop process by name(%v) fail with process is not active", name)
	}
	if proc.Status != PS_RUNNING {
		return util.Err("Runner stop process by name(%v) fail with process is not running", name)
	}
	proc.Cmd.Process.Kill()
	var used int64 = 0
	for proc.Status == PS_RUNNING {
		time.Sleep(100 * time.Millisecond)
		used += 100
		if used >= timeout {
			return util.Err("Runner stop process(%v) timeout", proc.Name)
		}
	}
	return nil
}

func (r *Runner) RunProc(p *Proc) error {
	r.Wg.Add(1)
	defer r.Wg.Done()
	log.I("Runner start process(%v) by cws(%v),exec(%v),args(%v),out(%v),err(%v),on(%v)",
		p.Name, p.Cws, p.Exec, p.Args, p.OutF, p.ErrF, p.On)
	var runner = exec.Command(p.Exec, util.ParseArgs(p.Args)...)
	runner.Dir = p.Cws
	var env = p.Envs
	if len(env) > 0 {
		runner.Env = append(os.Environ(), strings.Split(env, ",")...)
	}
	var out_w, err_w *smartio.TimeFlushWriter = nil, nil
	if len(p.OutF) > 0 {
		out_w = smartio.NewNamedWriter(p.Cws, p.OutF, 1024, r.WDelay)
		defer out_w.Stop()
		runner.Stdout = out_w
	}
	if len(p.ErrF) > 0 {
		err_w = smartio.NewNamedWriter(p.Cws, p.ErrF, 1024, r.WDelay)
		defer err_w.Stop()
		runner.Stderr = err_w
	}
	p.Cmd = runner
	p.Status = PS_RUNNING
	err := runner.Start()
	if err != nil {
		err = util.Err("Proc start process fail with %v", err)
		log.E("%v", err)
		p.Status = PS_START_ERR
		return err
	}
	err = runner.Wait()
	log.I("Runner is done for process(%v) with error(%v)", p.Name, err)
	p.Status, p.Cmd = PS_NOT_START, nil
	return err
}

func (r *Runner) Kill() {
	var sended = 0
	for _, proc := range r.ProcL {
		if proc.Cmd == nil || proc.Cmd.Process == nil {
			continue
		}
		proc.Cmd.Process.Kill()
		sended += 1
	}
	log.D("Runner send kill signal to %v process", sended)
}

func (r *Runner) SrvHTTP(hs *routing.HTTPSession) routing.HResult {
	var name, exec, token string
	var err = hs.ValidCheckVal(`
		name,R|S,L:0;
		exec,R|S,L:0;
		token,R|S,L:0;
		`, &name, &exec, &token)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	if r.Token[token] < 1 {
		err = util.Err("token(%v) not found", token)
		return hs.MsgResErr2(401, "arg-err", err)
	}
	switch exec {
	case "start":
		err = r.StartProc(name, 30000)
	case "stop":
		err = r.StopProc(name, 30000)
	case "restart":
		err = r.RestartProc(name, 30000)
	default:
		err = util.Err("command(%v) not found", exec)
	}
	if err == nil {
		return hs.MsgRes("OK")
	} else {
		return hs.MsgResErr2(2, "srv-err", err)
	}

}
