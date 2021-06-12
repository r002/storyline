package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	fmt.Println(">> Initializing testing:", filename)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestGetEnvVars(t *testing.T) {
	ev := GetEnvVars()
	fmt.Println(">> Test config service!!!!", ev)
}
