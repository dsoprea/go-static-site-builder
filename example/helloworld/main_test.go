package main

import (
    "io/ioutil"
    "os/exec"
    "testing"

    "github.com/dsoprea/go-logging"
)

func TestMain(t *testing.T) {
    tempPath, err := ioutil.TempDir("", "")
    log.PanicIf(err)

    cmd := exec.Command("go", "run", "main.go", "--output-path", tempPath)

    err = cmd.Run()
    log.PanicIf(err)
}
