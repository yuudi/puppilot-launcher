package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func entry() {
	if err := ensurePathExist("puppilot"); err != nil {
		fmt.Println(err.Error())
		return
	}
	nodePath, err := getNodePath()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// check for updates
	if err := checkForUpdate(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// start puppilot
	fmt.Println("Starting puppilot")
	jsPath, err := filepath.Abs(path.Join("puppilot", "puppilot.js"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cmd := exec.Command(nodePath, "--enable-source-maps", jsPath)
	cmd.Dir = "puppilot"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	entry()
	fmt.Println("Press enter to exit")
	fmt.Scanln()
}
