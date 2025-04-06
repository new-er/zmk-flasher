package main

import (
	"os"
	"runtime"
	"zmk-flasher/cmd"
	"zmk-flasher/platform"
)

func main() {

	switch runtime.GOOS {
	case "darwin":
		platform.Os = platform.DarwinOperations{}
	case "linux":
		platform.Os = platform.LinuxOsOperations{}
	default:
		println("OS not supported yet")
		os.Exit(1)
	}

	err := cmd.Execute()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
