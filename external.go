package skimmer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	EnvHttpBrowser   = "SKIM_HTTP_BROWSER"
	EnvGopherBrowser = "SKIM_GOPHER_BROWSER"
	EnvGemeniBrowser = "SKIM_GEMINI_BROWSER"
	EnvFtpBrowser    = "SKIM_FTP_BROWSER"
)

var clear map[string]func() //create a map for storing clear funcs

func SetupScreen(out io.Writer) {
	clear = make(map[string]func())
	for _, posix := range []string{"linux", "darwin"} {
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
	if ok {                          //if we defined a clear func for that platform:
		clrFn() //we execute it
	}
	// else we do nothing because we don't know how ..
}

func OpenInBrowser(in io.Reader, out io.Writer, eout io.Writer, link string) error {
	browser := ""
	switch {
	case strings.HasPrefix(link, "https:") || strings.HasPrefix(link, "http:"):
		// Open in the strimmer recommended web browser
		if browser = os.Getenv(EnvHttpBrowser); browser == "" {
			return fmt.Errorf("%s is not set", EnvHttpBrowser)
		}
	}
	if browser != "" {
		cmd := exec.Command(browser, link)
		cmd.Stdin = in
		cmd.Stdout = out
		cmd.Stderr = eout
		return cmd.Run()
	}
	return fmt.Errorf("I do not know how to open %q", link)
}
