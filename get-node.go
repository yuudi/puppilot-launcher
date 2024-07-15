package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/ulikunitz/xz"
)

func installNode_windows() error {
	downloadURL := "https://nodejs.org/dist/v22.4.1/node-v22.4.1-x64.msi"
	cmd := exec.Command("powershell", "Start-BitsTransfer", downloadURL, "C:\\temp\\node.msi")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("msiexec", "/i", "C:\\temp\\node.msi", "/quiet", "/qn", "/norestart")
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func checkNodeBinary_windows() (string, error) {
	nodePath := ".\\puppilot\\node.exe"
	_, err := os.Stat(nodePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return filepath.Abs(path.Join("puppilot", "node.exe"))
}

func getNodeBinary_windows() (string, error) {
	fmt.Println("Downloading nodejs")
	res, err := http.Get(nodejsUrlWindows)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	uncompressed, err := xz.NewReader(res.Body)
	if err != nil {
		return "", err
	}
	nodePath, err := filepath.Abs(path.Join("puppilot", "node.exe"))
	if err != nil {
		return "", err
	}
	nodeFile, err := os.Create(nodePath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(nodeFile, uncompressed); err != nil {
		return "", err
	}
	nodeFile.Close()
	return nodePath, nil
}

func getLocalNode() (string, error) {
	// check os
	os := runtime.GOOS
	switch os {
	case "windows":
		// return installNode_windows()
		nodePath, err := checkNodeBinary_windows()
		if err != nil {
			return "", err
		}
		if nodePath != "" {
			fmt.Println("Nodejs found puppilot directory")
			return nodePath, nil
		}
		path, err := getNodeBinary_windows()
		if err != nil {
			return "", err
		}
		return path, nil
	default:
		// not implemented yet
		return "", errors.New("OS not supported")
	}
}

func getNodePath() (string, error) {
	fmt.Println("Checking for node")
	nodePath, err := exec.LookPath("node")
	if err != nil {
		path, err := getLocalNode()
		if err != nil {
			return "", err
		}
		return path, nil
	}
	fmt.Println("Nodejs found in system path")
	return nodePath, nil
}
