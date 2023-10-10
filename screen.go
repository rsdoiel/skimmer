package skimmer

import (
    "io"
    "os/exec"
    "runtime"
)

var clear map[string]func() //create a map for storing clear funcs

func SetupScreen(out io.Writer) {
    clear = make(map[string]func())
	for _, posix := range []string{ "linux", "darwin" } {
    	clear[posix] = func() { 
        	cmd := exec.Command("clear") //Linux example, its tested
        	cmd.Stdout = out
        	cmd.Run()
    	}
	}
    clear["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = out
        cmd.Run()
    }
}

func ClearScreen() {
    clrFn, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
    if ok { //if we defined a clear func for that platform:
        clrFn()  //we execute it
    } 
	// else we do nothing because we don't know how ..
}

