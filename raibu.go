package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

var filesMap = map[string]time.Time{}

var (
	path     string
	build    string
	run      string
	progname string

	proc *exec.Cmd
	err  error
)

var (
	buildOptions *CommandOptions
	runOptions   *CommandOptions
)

func PrintLog(a ...any) {
	fmt.Println(progname, "::", fmt.Sprint(a...))
}

func ReadOutput(r io.ReadCloser) {
	buf := make([]byte, 1<<12)
	for {
		readed, err := r.Read(buf)
		if err == io.EOF {
			return
		}
		if err != nil {
			PrintLog("ReadOutput() =>", err)
			return
		}
		fmt.Print(string(buf[:readed]))
	}
}

func executeCommands() {
	for _, opt := range []*CommandOptions{buildOptions, runOptions} {
		if opt == nil {
			continue
		}
		PrintLog(opt.HeaderMsg)
		proc := exec.Command(opt.ExecName, opt.Args...)
		r, err := proc.StdoutPipe()
		if err != nil {
			PrintLog(opt.ErrorString(err))
			return
		}
		if err := proc.Start(); err != nil {
			PrintLog(opt.ErrorString(err))
			return
		}
		go ReadOutput(r)
		if err := proc.Wait(); err != nil {
			PrintLog(opt.ErrorString(err))
		}
	}
}

func main() {
	go executeCommands()
	for {
		scanFiles(path)
	}
}

func init() {
	progname = os.Args[0]
	flag.StringVar(&path, "path", ".", "Directory path to watch")
	flag.StringVar(&build, "build", "", "Build command to execute on files change")
	flag.StringVar(&run, "run", "", "Run command to execute on files change")
	flag.Parse()

	if run == "" {
		PrintLog("Run command is empty. Exit...")
		os.Exit(1)
	}
	buildOptions = MakeOptions(build, "\033[34mBuilding...\033[0m", "\033[31mError on build\033[0m, executable: %s, args: %#v, error: %s, output: %s")
	runOptions = MakeOptions(run, "\033[32mRunning...\033[0m", "\033[31mError on run\033[0m, executable: %s, args: %#v, error: %s, output: %s")
}
