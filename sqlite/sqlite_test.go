package sqlite

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bengtan/websqlview/native"
	"github.com/bengtan/websqlview/webviewex"
	"github.com/zserge/webview"
)

var (
	wd            string
	dummyExitCode int
)

func TestMain(m *testing.M) {
	wd, _ = os.Getwd()
	os.Exit(m.Run())
}

func TestJS(t *testing.T) {
	matches, _ := filepath.Glob("test/*.js")
	for _, filename := range matches {
		fmt.Printf("Testing: %s", filename)
		failure := testOneJS(filename)
		if failure == "" {
			fmt.Println(" - pass")
		} else {
			fmt.Println(" - FAIL")
			t.Errorf(failure)
			return
		}
	}
}

func testOneJS(filename string) (failure string) {
	failure = ""
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return err.Error()
	}

	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Testing: " + filename)
	w.SetSize(800, 600, webview.HintNone)
	w.Bind("pass", func() {
		w.Terminate()
	})
	w.Bind("fail", func(s string) {
		w.Terminate()
		failure = s
	})

	ex := webviewex.New(w)
	native.Init(ex, &dummyExitCode)

	Init(ex)
	defer Shutdown()

	// Override with sqlite.js
	sqliteJs, err := ioutil.ReadFile("sqlite.js")
	if err != nil {
		return err.Error()
	}
	w.Init(string(sqliteJs))

	w.Init(string(text))
	w.Navigate("file://" + wd + "/test/testHarness.html")
	w.Run()
	return
}
